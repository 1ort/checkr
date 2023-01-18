package provider

import (
	"log"

	"github.com/1ort/checkr/proxy"
)

type stringProvider struct {
	urls []string
}

func FromStrings(proxyUrls ...string) ProxyProvider {
	return &stringProvider{
		proxyUrls,
	}
}

func (p *stringProvider) Fetch(done chan int) <-chan proxy.Proxy {
	ch := make(chan proxy.Proxy)
	go func() {
		defer close(ch)

		for _, proxyurl := range p.urls {
			proxy, err := proxy.FromHostPort(proxyurl)
			if err != nil {
				log.Println(err)
				continue
			}
			select {
			case <-done:
				return
			default:
				ch <- proxy
			}
		}
	}()
	return ch
}
