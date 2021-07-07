package main

import (
	"github.com/vielendanke/go-retry/client"
	"github.com/vielendanke/go-retry/dialer"
	"net/http"
	"time"
)

func main() {
	f := dialer.SetKeepAliveProbe(3 * time.Second)
	s := dialer.SetTimeout(500 * time.Millisecond)
	dial := dialer.New(f, s)
	cli := client.New(client.WithDialerRetry(dial, 100 * time.Millisecond, 5))
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:9090/try", nil)
	cli.Do(req)
}