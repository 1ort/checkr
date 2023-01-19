package counter

import (
	"log"

	"github.com/1ort/checkr/proxy"
)

type ProxyCounter struct {
	Name   string
	Count  int
	Closed bool
}

func NewProxyCounter(name string) ProxyCounter {
	var c ProxyCounter
	c.Name = name
	return c
}

func (c *ProxyCounter) Run(in <-chan proxy.Proxy) <-chan proxy.Proxy {
	out := make(chan proxy.Proxy)
	go func() {
		for p := range in {
			out <- p
			c.Count++
		}
		c.Closed = true
		close(out)
	}()
	return out
}

func (c *ProxyCounter) Print() {
	log.Printf("%s: %v proxy", c.Name, c.Count)
}
