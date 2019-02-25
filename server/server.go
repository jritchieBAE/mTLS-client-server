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

	var server *http.Server
	var err error
	switch serverType {
	case None:
		server = mtls.NewServer(":8443")
		log.Fatal(server.ListenAndServe())
	case TLS:
		server, err = mtls.NewTlsServer(":8443", certPath, keyPath)
		log.Fatal(server.ListenAndServeTLS("", ""))
	case mTLS:
		server, err = mtls.NewMtlsServer(":8443", certPath, keyPath, caCertPath)
		log.Fatal(server.ListenAndServeTLS("", ""))
	}
	if err != nil {
		log.Fatal(err)
	}

}
