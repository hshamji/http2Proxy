package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"net/url"
)

func Oldmain() {
	c := http.Client{}
	//resp, err := c.Get("http://localhost:8080")
	r, _ := http.NewRequest("POST", "http://localhost:8080", nil)
	resp, err := c.Do(r)
	if err != nil {
		fmt.Printf("Error is: %+v\n", err)
	} else {
		fmt.Printf("Response: %+v\n", resp)
	}
}

func nomain() {
	u := url.URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "http://localhost:8080",
		Path:        "/tensorflow.serving.PredictionService/Predict",
		RawPath:     "",
		OmitHost:    false,
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	fmt.Printf("The extracted RequestURI is: %+v\n", u.RequestURI())

	u2, _ := http.NewRequest("POST", "http://127.0.0.1:8080", nil)
	fmt.Printf("The extracted RequestURI is: %+v\n", u2.URL.RequestURI())

}
func main() {

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

	///tensorflow.serving.PredictionService/Predict
	req2, _ := http.NewRequest("POST", "http://127.0.0.1:8080", nil)
	fmt.Printf("Complete request: %+v\n", req2)
	fmt.Printf("CompleteURL: %+v\n", req2.URL)
	req := http.Request{}
	req.URL = &url.URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "http://127.0.0.1:8080",
		Path:        "/tensorflow.serving.PredictionService/Predict",
		RawPath:     "",
		OmitHost:    false,
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	req2.URL.Path = "/tensorflow.serving.PredictionService/Predict"
	req = *req2 // Comment this out to use manual Request
	req.Proto = "HTTP/2.0"
	req.ProtoMajor = 2
	req.ProtoMinor = 0
	header := http.Header{}
	header.Set("Content-Type", "application/grpc")
	header.Set("Grpc-Accept-Encoding", "gzip")
	header.Set("Grpc-Encoding", "gzip")
	header.Add("User-Agent", "throttleProxy")
	//header.Add("Path", "/tensorflow.serving.PredictionService/Predict")
	req.Header = header
	//req.RequestURI = "/tensorflow.serving.PredictionService/Predict"
	//req.Host = "/tensorflow.serving.PredictionService/Predict"

	r := bytes.NewReader([]byte{1, 0, 0, 0, 85, 31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 226, 82, 224, 226, 205, 72, 204, 73, 139, 47, 200, 41, 45, 142, 47, 41, 207, 151, 226, 47, 78, 45, 42, 203, 204, 75, 143, 79, 73, 77, 75, 44, 205, 41, 17, 146, 228, 98, 172, 16, 18, 225, 96, 212, 18, 96, 96, 104, 176, 103, 96, 96, 112, 96, 96, 112, 112, 96, 96, 104, 112, 0, 4, 0, 0, 255, 255, 40, 53, 85, 182, 61, 0, 0, 0})
	rc := io.NopCloser(r)
	req.Body = rc

	fmt.Printf("The request is: %+v\n\n", req)

	resp, err := client.Do(&req)
	if err != nil {
		fmt.Printf("Error is %+v\n\n", err)
	} else {
		fmt.Printf("Response is %+v\n\n", resp)
	}
}

//RECEIVED Request: {Method:POST URL:/tensorflow.serving.PredictionService/Predict Proto:HTTP/2.0 ProtoMajor:2 ProtoMinor:0 Header:map[Content-Type:[application/grpc] Grpc-Accept-Encoding:[gzip] Grpc-Encoding:[gzip] Te:[trailers] User-Agent:[grpc-go/1.54.0]] Body:0xc00008c240 GetBody:<nil> ContentLength:-1 TransferEncoding:[] Close:false Host:localhost:8080 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:[::1]:59225 RequestURI:/tensorflow.serving.PredictionService/Predict TLS:<nil> Cancel:<nil> Response:<nil> ctx:0xc0000960a0}
//
//Body of Request: {buf:[1 0 0 0 85 31 139 8 0 0 0 0 0 0 255 226 82 224 226 205 72 204 73 139 47 200 41 45 142 47 41 207 151 226 47 78 45 42 203 204 75 143 79 73 77 75 44 205 41 17 146 228 98 172 16 18 225 96 212 18 96 96 104 176 103 96 96 112 96 96 112 112 96 96 104 112 0 4 0 0 255 255 40 53 85 182 61 0 0 0] off:0 lastRead:0}
