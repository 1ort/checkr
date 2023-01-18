package provider

import (
	"time"

	"github.com/1ort/checkr/proxy"
)

type repeatProvider struct {
	repeatTimes int
	sleepTime   time.Duration
	child       ProxyProvider
}

func Repeat(p ProxyProvider, times int, sleep time.Duration) ProxyProvider {
	return &repeatProvider{
		times,
		sleep,
		p,
	}
}

func (p *repeatProvider) Fetch(done chan int) <-chan proxy.Proxy {
	out := make(chan proxy.Proxy)
	go func() {
		defer close(out)
		for i := p.repeatTimes; i != 0; i-- {
			for prox := range p.child.Fetch(done) {
				out <- prox
			}
			select {
			case <-done:
				return
			default:
				time.Sleep(p.sleepTime)
			}
		}
	}()
	return out
}
