package provider

import (
	"sync"

	"github.com/1ort/checkr/proxy"
)

type providerGroup struct {
	providers []ProxyProvider
}

func merge[T any](cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	output := func(c <-chan T) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (p *providerGroup) Fetch(done chan int) <-chan proxy.Proxy {
	cs := make([]<-chan proxy.Proxy, len(p.providers))
	for i, p := range p.providers {
		cs[i] = p.Fetch(done)
	}
	return merge(cs...)
}

func Group(providers ...ProxyProvider) ProxyProvider {
	if len(providers) == 0 {
		return nil
	}
	return &providerGroup{
		providers,
	}
}
