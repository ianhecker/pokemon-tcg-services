package proxy

import (
	"net"
	"net/http"
	"time"
)

// For fast endpoints: quick connect, quick TLS, short header wait,
// healthy idle pool, and a sane ceiling on in-flight per host.
func NewTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,  // fast connect or bail
			KeepAlive: 30 * time.Second, // keep connections warm
		}).DialContext,

		// Reuse a single HTTP/2 TCP conn with many streams where possible.
		ForceAttemptHTTP2:   true,
		TLSHandshakeTimeout: 3 * time.Second,

		// Snappy APIs shouldn't stall; fail fast if headers don't arrive.
		// If you occasionally see timeouts on good networks, bump to 4–5s.
		ResponseHeaderTimeout: 2 * time.Second,

		// Mostly GETs; don't bother with 100-continue dance.
		ExpectContinueTimeout: 0,

		// Keep an idle pool so bursts don't pay new handshakes,
		// but drain idle conns reasonably quickly.
		IdleConnTimeout:     45 * time.Second,
		MaxIdleConns:        128,
		MaxIdleConnsPerHost: 64,

		// Hard ceiling to prevent stampede; tune to your parallelism.
		// Set to 0 if you enforce concurrency elsewhere.
		MaxConnsPerHost: 64,

		// Leave gzip on; small JSON gets smaller.
		DisableCompression: false,
	}
}
