package check

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/1ort/checkr/proxy"
)

// var ErrDeadProxy = errors.New("Proxy is dead")
// var ErrRequest = errors.New("Error requesting url")

var defaultReqAddr = "api.ip.sb/ip"

type checkerPool struct {
	in           <-chan proxy.Proxy
	out          chan<- proxy.Proxy
	reqAddr      string
	timeout      time.Duration
	workers      int
}

type ProxyChecker interface {
	Run(<-chan proxy.Proxy) <-chan proxy.Proxy
}

func NewCheckerPool(timeout time.Duration, workers int) ProxyChecker {
	return &checkerPool{
		nil,
		nil,
		defaultReqAddr,
		timeout,
		workers,
	}
}

func (c *checkerPool) Run(in <-chan proxy.Proxy) <-chan proxy.Proxy {
	var wg sync.WaitGroup
	out := make(chan proxy.Proxy)
	c.out = out
	c.in = in

	wg.Add(c.workers)
	for i := 0; i < c.workers; i++ {
		go func() {
			for p := range c.in {
				c.check(p)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		log.Println("checker: in channel closed")
		close(c.out)
	}()

	return out
}

func (c *checkerPool) check(p proxy.Proxy) {
	if p.Type == proxy.TypeUnknown {
		c.checkUnknown(p)
	} else {
		c.checkTyped(p)
	}
}

func (c *checkerPool) checkUnknown(p proxy.Proxy) {
	for _, pType := range proxy.PossibleTypes {
		pTyped := p
		pTyped.Type = pType
		c.checkTyped(pTyped)
	}
}

func (c *checkerPool) checkTyped(p proxy.Proxy) {
	schema := getSchema(p.Type)
	transport, err := p.Transport()
	if err != nil {
		p.IsAlive = false
		c.out <- p
		return //TODO: output err
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   c.timeout,
	}
	response, err := client.Get(schema + c.reqAddr)
	if err != nil {
		p.IsAlive = false
		c.out <- p
		return //TODO: output err
	}
	defer response.Body.Close()
	p.IsAlive = true
	c.out <- p

}

func getSchema(t proxy.ProxyType) string {
	switch t {
	case proxy.TypeHTTPS:
		return "https://"
	case proxy.TypeHTTP, proxy.TypeSOCKS4, proxy.TypeSOCKS5:
		return "http://"
	default:
		return ""
	}
}
