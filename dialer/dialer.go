package dialer

import (
	"net"
	"time"
)

type Option func(dialer *net.Dialer)

func SetTimeout(timeout time.Duration) Option {
	return func(dialer *net.Dialer) {
		dialer.Timeout = timeout
	}
}

func SetKeepAliveProbe(keepAlive time.Duration) Option {
	return func(dialer *net.Dialer) {
		dialer.KeepAlive = keepAlive
	}
}

func New(opts ...Option) *net.Dialer {
	d := &net.Dialer{}
	for _, v := range opts {
		v(d)
	}
	return d
}
