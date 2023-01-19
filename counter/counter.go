package counter

import (
	"log"
	"time"

	"github.com/1ort/checkr/proxy"
)

type proxyCounter struct {
	name   string
	closed bool
	count  int
}

func NewProxyCounter(name string, in <-chan proxy.Proxy, sleepTime time.Duration) <-chan proxy.Proxy {
	out := make(chan proxy.Proxy)
	var c proxyCounter
	c.name = name

	go func() {
		for p := range in {
			out <- p
			c.count++
		}
		c.closed = true
		close(out)
	}()

	go func() {
		for !c.closed {
			time.Sleep(sleepTime)
			c.print()
		}
	}()

	return out
}

func (c *proxyCounter) print() {
	log.Printf("%s: %v proxy", c.name, c.count)
}
