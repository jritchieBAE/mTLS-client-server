package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	mtls "../mtlsServer"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection accepted!")
	io.WriteString(w, "Hello world!\n")
}

type TLSType int

const (
	None TLSType = 0
	TLS  TLSType = 1
	mTLS TLSType = 2
)

var (
	certPath   = "../serverSigned.crt"
	keyPath    = "../server.key"
	caCertPath = "../root.crt"
)

func main() {

	http.HandleFunc("/hello", helloHandler)

	serverType := TLS

	var server *mtls.TlsServer
	var err error
	switch serverType {
	case None:
		server = mtls.NewUnsecureServer()
	case TLS:
		server, err = mtls.NewTlsServer(certPath, keyPath)
	case mTLS:
		server, err = mtls.NewMtlsServer(certPath, keyPath, caCertPath)
	}
	if err != nil {
		log.Fatal(err)
	}

	server.Listen(":8443")
}
