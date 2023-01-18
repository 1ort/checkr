package provider

import "github.com/1ort/checkr/proxy"

type ProxyProvider interface {
	Fetch(done chan int) <-chan proxy.Proxy
	// RunEndless() <-chan proxy.Proxy
}
