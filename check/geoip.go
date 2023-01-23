package check

import (
	"log"

	"github.com/1ort/checkr/mmdb"
	"github.com/1ort/checkr/proxy"
	"github.com/oschwald/maxminddb-golang"
)

type geoChecker struct {
	in     <-chan proxy.Proxy
	out    chan<- proxy.Proxy
	reader *maxminddb.Reader
}

func NewGeoChecker(reader *maxminddb.Reader) *geoChecker {
	var gc geoChecker
	gc.reader = reader
	return &gc
	// return &geoChecker{}
}

func (c *geoChecker) Run(in <-chan proxy.Proxy) <-chan proxy.Proxy {
	out := make(chan proxy.Proxy)
	c.out = out
	c.in = in

	go func() {
		for p := range c.in {
			// log.Printf("GeoChecker: Checking %v", p.Host)
			c.check(p)
		}
		log.Println("GeoChecker: in channel closed")
		close(c.out)
	}()

	return out
}

func (c *geoChecker) check(p proxy.Proxy) {
	geodata, err := mmdb.LookupGeo(c.reader, p.Host)
	if err != nil {
		return
	}
	p.Geo = geodata
	c.out <- p
}
