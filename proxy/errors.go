package proxy

import "errors"

var ErrInvalidHostPort = errors.New("invalid proxy host/port")
var ErrInvalidPort = errors.New("Port maximum value is 65535")
var ErrInvalidType = errors.New("Invalid proxy type")
