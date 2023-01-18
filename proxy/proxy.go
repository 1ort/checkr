package proxy

import (
	"net"
	"net/http"
	"strconv"
)

type ProxyType int

const (
	TypeUnknown ProxyType = iota
	TypeHTTP
	TypeHTTPS
	TypeSOCKS4
	TypeSOCKS4A
	TypeSOCKS5
	// TypeSMTP
)

var PossibleTypes = []ProxyType{TypeHTTP, TypeSOCKS4, TypeSOCKS5}

type Proxy struct {
	Host    string
	Port    int
	Type    ProxyType
	IsAlive bool
}

// type ProxyGeo struct {
// 	CountryISO string
// 	Country    string
// 	RegionIso  string
// 	Region     string
// 	City       string
// }

func (p *Proxy) HostPortString() string {
	return p.Host + ":" + strconv.Itoa(p.Port)
}

func (p *Proxy) UrlString() string {
	return p.Schema() + p.HostPortString()
}

func (p *Proxy) Schema() string {
	switch p.Type {
	case TypeHTTP:
		return "http://"
	case TypeHTTPS:
		return "https://"
	case TypeSOCKS4:
		return "socks4://"
	case TypeSOCKS4A:
		return "socks4a://"
	case TypeSOCKS5:
		return "socks5://"
	default:
		return ""
	}
}

func FromHostPort(hostport string) (Proxy, error) {
	var (
		p       Proxy
		err     error
		portStr string
	)
	p.Host, portStr, err = net.SplitHostPort(hostport)
	if err != nil {
		return p, ErrInvalidHostPort
	}
	p.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return p, ErrInvalidHostPort
	}
	if p.Port > 65535 {
		return p, ErrInvalidPort
	}
	return p, nil
}

func (p *Proxy) Transport() (*http.Transport, error) {
	switch p.Type {
	case TypeHTTP:
		return p.httpTransport()
	case TypeHTTPS:
		return p.httpsTransport()
	case TypeSOCKS4, TypeSOCKS4A, TypeSOCKS5:
		return p.socksTransport()
	default:
		return nil, ErrInvalidType
	}
}
