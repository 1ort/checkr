package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/1ort/checkr/check"
	"github.com/1ort/checkr/provider"
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
	inputFiles arrayFlags
	inputUrls  arrayFlags
	// outputFile string
	noCheck bool
	// silent     bool
	limit     int
	workers   int
	proxyType string
	timeout   int
)

var collected int

func main() {
	flag.Var(&inputFiles, "file", "Proxy list file. You can specify multiple (--file 1.txt --file 2.txt ... --file n.txt)")
	flag.Var(&inputUrls, "url", "Proxy list url You can specify multiple")
	// flag.StringVar(&outputFile, "o", "", "Output file")
	flag.BoolVar(&noCheck, "nocheck", false, "Do not check proxy")
	flag.IntVar(&limit, "limit", 0, "Limit proxies. 0 = no limit")
	flag.StringVar(&proxyType, "type", "all", "Type of proxy needed. [all/http/socks4/socks5] Doesn't work with --nocheck")
	// flag.BoolVar(&silent, "silent", false, "Enable silent mode")
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
	if !noCheck {
		checker := check.NewCheckerPool(
			time.Duration(timeout)*time.Second,
			false,
			workers,
		)
		out = checker.Run(out)
	}
	for pr := range out {
		if !pr.IsAlive {
			continue
		}
		collected++
		fmt.Println(pr.UrlString())
		if collected >= limit && limit != 0 {
			fmt.Printf("Limit reached\n")
			break
		}
	}
	fmt.Printf("Execution completed. %v proxy found\n", collected)
}
