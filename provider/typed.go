package provider

import "github.com/1ort/checkr/proxy"

type typedProvider struct {
	child     ProxyProvider
	proxyType proxy.ProxyType
}

func Socks5(providers ...ProxyProvider) ProxyProvider {
	return WithType(proxy.TypeSOCKS5, providers...)
}

func Socks4(providers ...ProxyProvider) ProxyProvider {
	return WithType(proxy.TypeSOCKS4, providers...)
}

func HTTP(providers ...ProxyProvider) ProxyProvider {
	return WithType(proxy.TypeHTTP, providers...)
}

func HTTPS(providers ...ProxyProvider) ProxyProvider {
	return WithType(proxy.TypeHTTPS, providers...)
}

func WithType(proxyType proxy.ProxyType, providers ...ProxyProvider) ProxyProvider {
	if len(providers) == 0 {
		return nil
	}
	if len(providers) == 1 {
		return &typedProvider{
			providers[0],
			proxyType,
		}
	}
	return &typedProvider{
		&providerGroup{
			providers,
		},
		proxyType,
	}
}

func (p *typedProvider) Fetch(done chan int) <-chan proxy.Proxy {
	out := make(chan proxy.Proxy)
	go func() {
		defer close(out)
		for prox := range p.child.Fetch(done) {
			prox.Type = p.proxyType
			out <- prox
		}
	}()
	return out
}
