package networking

import (
	"net"
	"net/http"
	"time"
)

func NewTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,

		// Let Go negotiate HTTP/2 for better reuse/multiplexing.
		ForceAttemptHTTP2:   true,
		TLSHandshakeTimeout: 10 * time.Second,

		// Important: let your per-request context control total time.
		// If the API can take ~1m to respond, don't cap header wait here.
		ResponseHeaderTimeout: 0,

		// You’re not sending large bodies (mostly GETs), so 100-continue isn’t needed.
		ExpectContinueTimeout: 0,

		// Pooling: generous but not excessive for a rate-limited API.
		IdleConnTimeout:     120 * time.Second,
		MaxIdleConns:        64,
		MaxIdleConnsPerHost: 32,

		// Leave unlimited; enforce rate/parallelism in your own limiter.
		// (If you want a hard ceiling on in-flight to this host, set 4–8 here.)
		MaxConnsPerHost: 0,

		// Keep gzip on for smaller JSON payloads.
		DisableCompression: false,
	}
}
