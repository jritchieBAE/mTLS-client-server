package main

import (
	"fmt"
	"io"
	"net/http"

	mtls "../pkg"
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

func main() {

	http.HandleFunc("/hello", helloHandler)

	serverType := mTLS

	var server *mtls.TlsServer
	switch serverType {
	case None:
		server = mtls.NewUnsecureServer()
	case TLS:
		server = mtls.NewTlsServer("../cert.pem", "../key.pem")
	case mTLS:
		server = mtls.NewMtlsServer("../cert.pem", "../key.pem")
	}

	server.Listen(":8443")
}
