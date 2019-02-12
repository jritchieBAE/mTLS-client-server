package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

	switch serverType {
	case None:
		UnsecureServer()
	case TLS:
		TlsServer()
	case mTLS:
		MtlsServer()
	}
}

func MtlsServer() {

	caCert, err := ioutil.ReadFile("../cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	log.Fatal(server.ListenAndServeTLS("../cert.pem", "../key.pem"))
}

func TlsServer() {

	log.Fatal(http.ListenAndServeTLS(":8443", "../cert.pem", "../key.pem", nil))
}

func UnsecureServer() {

	log.Fatal(http.ListenAndServe(":8443", nil))
}
