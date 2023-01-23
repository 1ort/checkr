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
	Geo     *GeoInfo
}

type GeoInfo struct {
	Continent struct {
		Code  string            `maxminddb:"code" json:"code"`
		Names map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"continent" json:"continent"`
	Country struct {
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union" json:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code" json:"iso_code"`
		Names             map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"country" json:"country"`
	Subdivisions []struct {
		IsoCode string            `maxminddb:"iso_code" json:"iso_code"`
		Names   map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"subdivisions" json:"subdivisions"`
	City struct {
		Names map[string]string `maxminddb:"names" json:"names"`
	} `maxminddb:"city" json:"city"`
	Postal struct {
		Code string `maxminddb:"code" json:"code"`
	} `maxminddb:"postal" json:"postal"`
	Location struct {
		Latitude  float64 `maxminddb:"latitude" json:"latitude"`
		Longitude float64 `maxminddb:"longitude" json:"longitude"`
		TimeZone  string  `maxminddb:"time_zone" json:"time_zone"`
	} `maxminddb:"location" json:"location"`
}

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
	// case TypeHTTPS:
	// 	return p.httpsTransport()
	case TypeSOCKS4, TypeSOCKS4A, TypeSOCKS5:
		return p.socksTransport()
	default:
		return nil, ErrInvalidType
	}
}
