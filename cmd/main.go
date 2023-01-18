package main

import (
	"fmt"
	"log"
	"time"

	"github.com/1ort/checkr/check"
	"github.com/1ort/checkr/provider"
)

func main() {
	p := provider.FromUrls(
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt",
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt",
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
		"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/proxy.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies.txt",
	)
	checker := check.NewCheckerPool(10*time.Second, false, 500)
	done := make(chan int)
	out := p.Fetch(done)
	out = checker.Run(out)
	for proxy := range out {
		if proxy.IsAlive {
			fmt.Println(proxy.UrlString())
		} else {
			fmt.Println("dead(")
		}
	}
	log.Println("out channel closed")
}
