# Checkr - Blazing fast proxy checker

![изображение](https://user-images.githubusercontent.com/83316072/213428651-a98dd5d9-3000-49c0-a637-d3e3b66021a4.png)

# Installation

from source:
To install from source, you need git and GO
```
git clone https://github.com/1ort/checkr.git
cd checkr/cmd/
go build . -o checkr
```

# Usage
`checkr -file unchecked.txt -url https://www.proxy-list.download/api/v1/get?type=http -type http -o checked.txt`

If proxy sources are not specified, checkr will collect proxies from public lists

```
  -file value
        Proxy list file. You can specify multiple (-file 1.txt -file 2.txt ... -file n.txt)
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

## Notes
If you do not specify the proxy format, or select all, the check speed will be much slower, since each proxy will be checked for compliance with each protocol. Since the server can support different protocols on the same port, these proxies will be recorded separately. For example:
```
socks5://151.151.151.151:123
socks4://151.151.151.151:123
http://151.151.151.151:123
```

# TODO

1) Support proxy with authentication
2) Checking proxy by geo
3) built-in server with rotation
4) saving the proxy to the database
