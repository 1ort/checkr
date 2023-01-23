package mmdb

import (
    "net"

    "github.com/1ort/checkr/proxy"
    "github.com/oschwald/maxminddb-golang"
)

func LookupGeo(db *maxminddb.Reader, ipStr string) (*proxy.GeoInfo, error) {
	var geoInfo proxy.GeoInfo
	ip := net.ParseIP(ipStr)

	err := db.Lookup(ip, &geoInfo)
	if err != nil {
		return nil, err
	}
	return &geoInfo, nil
}
