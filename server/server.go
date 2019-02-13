package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"../mtls"
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
	var err error
	switch serverType {
	case None:
		server = mtls.NewUnsecureServer()
	case TLS:
		server, err = mtls.NewTlsServer("../cert.pem", "../key.pem")
	case mTLS:
		server, err = mtls.NewMtlsServer("../cert.pem", "../key.pem")
	}

	if err != nil {
		log.Fatal(err)
	}

	server.Listen(":8443")
}
