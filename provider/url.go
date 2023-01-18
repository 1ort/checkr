package provider

import (
	"bufio"
	"log"
	"net/http"

	"github.com/1ort/checkr/proxy"
)

func FromUrls(urls ...string) ProxyProvider {
	if len(urls) == 0 {
		return nil
	}
	ps := make([]ProxyProvider, len(urls))
	for i, url := range urls {
		ps[i] = &urlProvider{
			url,
		}
	}
	return Group(ps...)
}

type urlProvider struct {
	url string
}

func (p *urlProvider) Fetch(done chan int) <-chan proxy.Proxy {
	ch := make(chan proxy.Proxy)
	go func() {
		defer close(ch)
		resp, err := http.Get(p.url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			proxy, err := proxy.FromHostPort(scanner.Text())
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
