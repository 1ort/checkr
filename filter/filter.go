package filter

import "github.com/1ort/checkr/proxy"

type ProxyFilterFunc func(p proxy.Proxy) bool

type proxyFilter struct {
	filters []ProxyFilterFunc
}

func NewProxyFilter(in <-chan proxy.Proxy, funcs ...ProxyFilterFunc) <-chan proxy.Proxy {
	if len(funcs) == 0 {
		return in
	}
	out := make(chan proxy.Proxy)
	f := proxyFilter{
		funcs,
	}
	go func() {
		for p := range in {
			if f.check(p) {
				out <- p
			}
		}
		close(out)
	}()
	return out
}

func Alive(p proxy.Proxy) bool {
	return p.IsAlive
}

func StrictCountryFilter(countryCode string) ProxyFilterFunc {
	return func(p proxy.Proxy) bool {
		if p.Geo == nil {
			return false
		}
		return p.Geo.Country.IsoCode == countryCode
	}
}

func OneOfCountryFilter(countryCodes ...string) ProxyFilterFunc {
	return func(p proxy.Proxy) bool {
		if p.Geo == nil {
			return false
		}
		for _, cc := range countryCodes {
			if p.Geo.Country.IsoCode == cc {
				return true
			}
		}
		return false
	}
}

func (f *proxyFilter) check(p proxy.Proxy) bool {
	for _, filterFunc := range f.filters {
		if !filterFunc(p) {
			return false
		}
	}
	return true
}
