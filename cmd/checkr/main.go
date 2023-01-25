package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1ort/checkr/check"
	"github.com/1ort/checkr/counter"
	"github.com/1ort/checkr/filter"
	"github.com/1ort/checkr/mmdb"
	"github.com/1ort/checkr/provider"
	"github.com/1ort/checkr/proxy"
	"github.com/oschwald/maxminddb-golang"
)

var defaultUrls = []string{
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt",
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt",
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
	"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/proxy.txt",
	"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies.txt",
	"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks4.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt",
	"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
	"https://www.proxy-list.download/api/v1/get?type=http",
	"https://www.proxy-list.download/api/v1/get?type=https",
	"https://www.proxy-list.download/api/v1/get?type=socks4",
	"https://www.proxy-list.download/api/v1/get?type=socks5",
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	inputFiles     arrayFlags
	inputUrls      arrayFlags
	countryFilters arrayFlags
	mmKey          string
	outputFile     string
	noCheck        bool
	silent         bool
	workers        int
	proxyType      string
	timeout        int
)

func main() {
	flag.Var(&inputFiles, "file", "Proxy list file. You can specify multiple (-file 1.txt -file 2.txt ... -file n.txt)")
	flag.Var(&inputUrls, "url", "Proxy list url. You can specify multiple (-url a.com -url b.com ... --url n.com)")

	flag.Var(&countryFilters, "country", "Filter proxies by country IsoCode. You can specify multiple")
	flag.StringVar(&mmKey, "mmKey", "", "MaxMind License key to download and update GeoLite2-City db.")

	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.BoolVar(&noCheck, "nocheck", false, "Do not check proxy")
	flag.StringVar(&proxyType, "type", "all", "Type of proxy needed. [all/http/socks4/socks5] Doesn't work with --nocheck")
	flag.BoolVar(&silent, "silent", false, "Enable silent mode")
	flag.IntVar(&timeout, "timeout", 10, "Check timeout in seconds")
	flag.IntVar(&workers, "workers", 500, "Number of parallel checkers")
	flag.Parse()
	p := provider.Group(
		provider.FromFiles(inputFiles...),
		provider.FromUrls(inputUrls...),
	)
	if p == nil {
		p = provider.FromUrls(defaultUrls...)
	}
	if !noCheck {
		switch proxyType {
		case "socks5":
			p = provider.Socks5(p)
		case "socks4":
			p = provider.Socks4(p)
		case "http":
			p = provider.HTTP(p)
		case "all":
		default:
			panic("Incorrect proxy type")
		}
	}

	done := make(chan int)
	defer close(done)
	out := p.Fetch(done)

	counters := make([]*counter.ProxyCounter, 0)
	foundCounter := counter.NewProxyCounter("Found")
	out = foundCounter.Run(out)
	counters = append(counters, &foundCounter)
	// out = counters[len(counters)-1].Run(out)

	if len(countryFilters) > 0 {
		var (
			r   *maxminddb.Reader
			err error
		)
		exPath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exDir := filepath.Dir(exPath)
		mmdbDir := filepath.Join(exDir, "mmdb/")
		if mmKey != "" {
			downloader := mmdb.NewDownloader(mmKey, mmdbDir)
			downloader.Verbose = !silent
			r, err = downloader.Latest("GeoLite2-City")
			if err != nil {
				panic(err)
			}
		} else {
			r, err = mmdb.OpenLocal(mmdbDir, "GeoLite2-City")
			if err != nil {
				panic(err)
			}
		}
		geoChecker := check.NewGeoChecker(r)
		out = geoChecker.Run(out)

		out = filter.NewProxyFilter(out, filter.OneOfCountryFilter(countryFilters...))

		countryCounter := counter.NewProxyCounter("Desired country")
		out = countryCounter.Run(out)
		counters = append(counters, &countryCounter)
	}

	if !noCheck {
		checker := check.NewCheckerPool(
			time.Duration(timeout)*time.Second,
			workers,
		)
		out = checker.Run(out)
		out = filter.NewProxyFilter(out, filter.Alive)

		aliveCounter := counter.NewProxyCounter("Alive")
		out = aliveCounter.Run(out)
		counters = append(counters, &aliveCounter)

	}

	if !silent {
		go func() {
			for {
				time.Sleep(5 * time.Second)
				for _, cr := range counters {
					cr.Print()
				}
			}
		}()
	}

	if outputFile != "" {
		writeToFile(outputFile, out)
	} else {
		for p := range out {
			fmt.Println(p.UrlString())
		}
	}
	fmt.Printf("Execution completed. %v proxy found\n", counters[len(counters)-1].Count)
}

func writeToFile(filename string, in <-chan proxy.Proxy) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	for p := range in {
		_, err := fmt.Fprintln(f, p.UrlString())
		if err != nil {
			log.Panic(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
