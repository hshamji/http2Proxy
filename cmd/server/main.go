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

// const PORT = "78"
const PORT = "8080"

func main() {
	h2s := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("RECEIVED Request: %+v\n\n", *r)

		bodyReader := r.Body
		//r.GetBody()

		//buf := new(strings.Builder)
		buf := bytes.Buffer{}
		io.Copy(&buf, bodyReader)
		fmt.Printf("Body of Request: %+v\n\n", buf)
		//fmt.Fprintf(w, "Hello, %v, http: %v", r.URL.Path, r.TLS == nil)
		//http.Error(w, "Some error", http.StatusGatewayTimeout)

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

		r.Host = fmt.Sprintf("http://localhost:%s", PORT)
		resp, err := client.Do(r)
		if err != nil {
			fmt.Printf("Error is: %+v\n", err)
		} else {
			fmt.Printf("Response is: %+v\n\n", resp)
		}

	})
	//a.Write([]byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 226, 18, 228, 226, 40, 74, 45, 46, 200, 207, 43, 78, 21, 98, 181, 98, 102, 100, 98, 6, 4, 0, 0, 255, 255, 83, 46, 113, 12, 19, 0, 0, 0})
	//a.Write([]byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 226, 82, 224, 226, 205, 72, 204, 73, 139, 47, 200, 41, 45, 142, 47, 41, 207, 151, 226, 47, 78, 45, 42, 203, 204, 75, 143, 79, 73, 77, 75, 44, 205, 41, 17, 146, 228, 98, 172, 16, 18, 225, 96, 212, 18, 96, 96, 104, 176, 103, 96, 96, 112, 96, 96, 112, 112, 96, 96, 104, 112, 0, 4, 0, 0, 255, 255, 40, 53, 85, 182, 61, 0, 0, 0})

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

//Listening [0.0.0.0:8080]...
//RECEIVED Request: {Method:POST URL:/tensorflow.serving.PredictionService/Predict Proto:HTTP/2.0 ProtoMajor:2 ProtoMinor:0 Header:map[Content-Type:[application/grpc] Grpc-Accept-Encoding:[gzip] Te:[trailers] User-Agent:[grpc-go/1.54.0]] Body:0xc00008c240 GetBody:<nil> ContentLength:-1 TransferEncoding:[] Close:false Host:localhost:8080 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:[::1]:59548 RequestURI:/tensorflow.serving.PredictionService/Predict TLS:<nil> Cancel:<nil> Response:<nil> ctx:0xc0000960a0}
//
//Body of Request: {buf:[0 0 0 0 61 10 32 10 13 104 97 108 102 95 112 108 117 115 95 116 119 111 26 15 115 101 114 118 105 110 103 95 100 101 102 97 117 108 116 18 25 10 1 120 18 20 8 1 42 16 0 0 128 63 0 0 0 64 0 0 64 64 0 0 128 64] off:0 lastRead:0}
