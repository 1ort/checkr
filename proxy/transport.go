package proxy

import (
	"net/http"
	"net/url"

	"h12.io/socks"
)

func (p *Proxy) socksTransport() (*http.Transport, error) {
	dial := socks.Dial(p.UrlString())
	transport := &http.Transport{
		Dial: dial,
	}
	return transport, nil
}

func (p *Proxy) httpTransport() (*http.Transport, error) {
	proxyURL, err := url.Parse(p.UrlString())
	if err != nil {
		return nil, ErrInvalidHostPort
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return transport, nil
}

// func (p *Proxy) httpsTransport() (*http.Transport, error) {
// 	proxyURL, err := url.Parse(p.UrlString())
// 	if err != nil {
// 		return nil, ErrInvalidHostPort
// 	}
// 	transport := &http.Transport{
// 		Proxy: http.ProxyURL(proxyURL),
// 		// TLSClientConfig: &tls.Config{},
// 	}
// 	return transport, nil
// }
