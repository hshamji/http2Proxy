package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	h := newReverseProxy()
	h = allowh2c(h)
	http.ListenAndServe(":78", h)
}

func allowh2c(next http.Handler) http.Handler {
	h2server := &http2.Server{IdleTimeout: time.Second * 60}
	return h2c.NewHandler(next, h2server)
}

const targetHost = "localhost:8080"

func newReverseProxy() http.Handler {
	director := func(req *http.Request) {
		origHost := req.Host
		log.Printf("[director] host=%s url=%s", origHost, req.URL)
		req.URL.Scheme = "https"
		req.URL.Host = targetHost
		req.Host = targetHost
		req.Header.Set("host", targetHost)

		log.Printf("[director] rewrote host=%s to=%q", origHost, req.URL)
	}

	transport := loggingTransport{next: &http2.Transport{
		// So http2.Transport doesn't complain the URL scheme isn't 'https'
		AllowHTTP: true,
		// Pretend we are dialing a TLS endpoint.
		// Note, we ignore the passed tls.Config
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}}
	v := &httputil.ReverseProxy{
		Director:      director,
		Transport:     transport,
		FlushInterval: -1, // do not buffer streaming responses
	}
	return v
}

type loggingTransport struct {
	next http.RoundTripper
}

func (l loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	log.Printf("[proxy] start: %s url=%s hdrs=%d trailers=%d",
		req.Method, req.URL, len(req.Header), len(req.Trailer))
	for k, v := range req.Header {
		log.Printf("	> HDR %s=%#v", k, v)
	}
	for k, v := range req.Trailer {
		log.Printf("	> TRAILER %s=%#v", k, v)
	}

	defer func() {
		log.Printf("[proxy]   end: %s url=%s took=%s",
			req.Method, req.URL, time.Since(start).Truncate(time.Millisecond))
	}()

	resp, err := l.next.RoundTrip(req)

	if resp != nil {
		log.Printf("[proxy]   resp: code=%d hdrs=%d trailers=%d",
			resp.StatusCode, len(resp.Header), len(resp.Trailer))
		for k, v := range resp.Header {
			log.Printf("	< HDR %s=%#v", k, v)
		}
		for k, v := range resp.Trailer {
			log.Printf("	< TRAILER %s=%#v", k, v)
		}
	}
	return resp, err
}
