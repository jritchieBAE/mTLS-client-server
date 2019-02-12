package main

import (
	"io"
	"log"
	"net/http"

	mtls "../pkg"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!\n")
}

type TLSType int

const (
	None TLSType = 0
	TLS  TLSType = 1
	mTLS TLSType = 2
)

func main() {

	http.HandleFunc("/hello", helloHandler)

	serverType := TLS

	server := mtls.NewTLSServer()

	switch serverType {
	case None:
		{
			log.Fatal(server.ListenNoTLS(":8443"))
		}

	case TLS:
		{
			server, err := (mtls.NewTLSServer().WithCertificates("../cert.pem", "../key.pem"))
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(server.ListenTLS(":8443"))
		}

	case mTLS:
		{
			server, err := (mtls.NewTLSServer().WithCertificates("../cert.pem", "../key.pem"))
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(server.ListenMTLS(":8443"))
		}
	}
}
