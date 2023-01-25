## Attention! The project is at the initial stage of development, so there may be changes that break backward compatibility

# Checkr - Blazing fast proxy checker

![изображение](https://user-images.githubusercontent.com/83316072/213428651-a98dd5d9-3000-49c0-a637-d3e3b66021a4.png)

# Installation

from source:
To install from source, you need git and GO
```
git clone https://github.com/1ort/checkr.git
cd checkr/cmd/
go build ./checkr -o checkr
```

via go install:
```
go install github.com/1ort/checkr
```
or download precompiled binary:
https://github.com/1ort/checkr/releases

# Usage
`checkr -file unchecked.txt -url https://www.proxy-list.download/api/v1/get?type=http -type http -o checked.txt -country USA`

If proxy sources are not specified, checkr will collect proxies from public lists

```
  -country value
        Filter proxies by country IsoCode. You can specify multiple
  -file value
        Proxy list file. You can specify multiple (-file 1.txt -file 2.txt ... -file n.txt)
  -mmKey string
        MaxMind License key to download and update GeoLite2-City db.
  -nocheck
        Do not check proxy
  -o string
        Output file
  -silent
        Enable silent mode
  -timeout int
        Check timeout in seconds (default 10)
  -type string
        Type of proxy needed. [all/http/socks4/socks5] Doesn't work with --nocheck (default "all")
  -url value
        Proxy list url. You can specify multiple (-url a.com -url b.com ... --url n.com)
  -workers int
        Number of parallel checkers (default 500)
```

## Proxy formats
Input format: `host:port`

Output format: `schema://host:port` where `schema` can be `socks5/socks4/http`

## IPGeo Data

The application takes IP geolocation data from the GeoLite2-City database provided by MaxMind. To use an up-to-date database, you need to specify the License Key as the --mmKey flag.

`checkr -country US -mmKey xxxxxxxxxxxxxxxx`

You can get the license key after creating an account for free at this link: https://www.maxmind.com/en/accounts/current/license-key

You do not need to specify the key as a flag each time, since the database file is saved to the mmdb folder next to the executable file.
Each time you use -mmKey, the local database is checked for relevance and the newest version is downloaded.

It is recommended to update the geo base at least once a week

## Notes
If you do not specify the proxy format, or select all, the check speed will be much slower, since each proxy will be checked for compliance with each protocol. Since the server can support different protocols on the same port, these proxies will be recorded separately. For example:
```
socks5://151.151.151.151:123
socks4://151.151.151.151:123
http://151.151.151.151:123
```

# TODO

- [ ] Support proxy with authentication
- [x] Checking proxy by geo
- [ ] built-in server with rotation
- [ ] saving the proxy to the database
