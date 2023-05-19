package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io"
	"net"
	"net/http"
)

const PORT = "78"

//const PORT = "8080"

func main() {
	h2s := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("RECEIVED Request: %+v\n\n", *r)

		buf := bytes.Buffer{}
		io.Copy(&buf, r.Body)
		r.Body.Close()
		//fmt.Printf("Body of Request: %+v\n\n", buf)

		client := http.Client{
			Transport: &http2.Transport{
				// So http2.Transport doesn't complain the URL scheme isn't 'https'
				AllowHTTP: true,
				// Pretend we are dialing a TLS endpoint.
				// Note, we ignore the passed tls.Config
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
		}

		println("Creating new request")
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8080", nil)
		req.URL.Path = "/tensorflow.serving.PredictionService/Predict"
		//req.URL.Path = "/tensorflow.serving.PredictionService/Predict"
		//req.URL.Host = "http://127.0.0.1:8080"
		////req = *req2 // Comment this out to use manual Request
		req.Proto = "HTTP/2.0"
		req.ProtoMajor = 2
		req.ProtoMinor = 0

		header := http.Header{}
		//println("Copying headers")
		//copyHeader(header, r.Header)

		header.Set("Content-Type", "application/grpc")
		header.Set("Grpc-Accept-Encoding", "gzip")
		header.Set("Grpc-Encoding", "gzip")
		header.Add("User-Agent", "throttleProxy")

		req.Header = header
		println("Creating body")
		req.Body = io.NopCloser(&buf)
		//io.Copy(req.Body, r.Body)

		//header.Set("Content-Type", "application/grpc")
		//header.Set("Grpc-Accept-Encoding", "gzip")
		//header.Set("Grpc-Encoding", "gzip")
		//header.Add("User-Agent", "throttleProxy")
		////header.Add("Path", "/tensorflow.serving.PredictionService/Predict")
		//req.Header = header
		////req.RequestURI = "/tensorflow.serving.PredictionService/Predict"
		////req.Host = "/tensorflow.serving.PredictionService/Predict"
		//
		//read := bytes.NewReader([]byte{0, 0, 0, 0, 61, 10, 32, 10, 13, 104, 97, 108, 102, 95, 112, 108, 117, 115, 95, 116, 119, 111, 26, 15, 115, 101, 114, 118, 105, 110, 103, 95, 100, 101, 102, 97, 117, 108, 116, 18, 25, 10, 1, 120, 18, 20, 8, 1, 42, 16, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 0, 0, 128, 64})
		//rc := io.NopCloser(read)
		//req.Body = rc
		//r.URL.Path = "/tensorflow.serving.PredictionService/Predict"
		//r.URL.Host = "http://127.0.0.1:8080"
		//revProxy := httputil.NewSingleHostReverseProxy(req.URL)
		//rp := NewProxy

		//t := &http2.Transport{
		//	// So http2.Transport doesn't complain the URL scheme isn't 'https'
		//	AllowHTTP: true,
		//	// Pretend we are dialing a TLS endpoint.
		//	// Note, we ignore the passed tls.Config
		//	DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
		//		return net.Dial(network, addr)
		//	},
		//}
		//revProxy.Transport = t
		//revProxy.ServeHTTP(w, r)

		println("About to send request")
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("Error is: %+v\n", err)
		} else {
			fmt.Printf("Response is: %+v\n\n", resp)
			//respBuf := bytes.Buffer{}
			//io.Copy(&respBuf, resp.Body)
			//fmt.Printf("ResponseBody is:%+v\n", respBuf)
			//	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			//	w.Header().Set("Grpc-Encoding", resp.Header.Get("Grpc-Encoding"))
			copyHeader(w.Header(), resp.Header)
			//copyHeader(w.Header(), resp.Trailer)
			//w.Header().Set(http2.TrailerPrefix, resp.Trailer)
			//w.Header().Set(http2.TrailerPrefix, "end")
			//defer func() {
			//	w.Header().Set("end", "this")
			//}()
			//resp.Body.Close()
			io.Copy(w, resp.Body)
			resp.Body.Close()
			//w.Write(io.EOF)
		}

	})

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", PORT),
		Handler: h2c.NewHandler(handler, h2s),
	}

	fmt.Printf("Listening [0.0.0.0:%s]...\n", PORT)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		panic("Error serving")
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
