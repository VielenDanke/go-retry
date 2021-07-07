package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Option func(cli *http.Client)

func WithDialerRetry(dialer *net.Dialer, retryPause time.Duration, numRetries uint) Option {
	return func(cli *http.Client) {
		cli.Transport.(*http.Transport).DialContext = dialContextWithRetry(dialer, retryPause, numRetries)
	}
}

func WithDialerRetryTLS(dialer *net.Dialer, retryPause time.Duration, numRetries uint) Option {
	return func(cli *http.Client) {
		cli.Transport.(*http.Transport).DialTLSContext = dialContextWithRetry(dialer, retryPause, numRetries)
	}
}

func WithProxy(proxy *url.URL) Option {
	return func(cli *http.Client) {
		cli.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(cli *http.Client) {
		cli.Transport.(*http.Transport).TLSClientConfig = cfg
	}
}

func WithDefaultDialer() Option {
	return func(cli *http.Client) {
		cli.Transport.(*http.Transport).DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}
}

func NewWithTLS(opts ...Option) *http.Client {
	tr := &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	cli := &http.Client{Transport: tr}
	for _, v := range opts {
		v(cli)
	}
	return cli
}

func New(opts ...Option) *http.Client {
	tr := &http.Transport{
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 1 * time.Second,
	}
	cli := &http.Client{Transport: tr}
	for _, v := range opts {
		v(cli)
	}
	return cli
}

func dialContextWithRetry(dialer *net.Dialer, retryPause time.Duration, numRetries uint) func(ctx context.Context, network string, addr string) (net.Conn, error) {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		var lastError error
		for i := uint(0); i < numRetries; i++ {
			conn, connErr := dialer.DialContext(ctx, network, addr)
			if connErr == nil {
				return conn, connErr
			}
			lastError = connErr
			t := time.NewTimer(retryPause)
			select {
			case <-ctx.Done():
				t.Stop()
				return nil, fmt.Errorf("context timed out: %s", ctx.Err())
			case <-t.C:
				t.Stop()
			}
		}
		return nil, lastError
	}
}
