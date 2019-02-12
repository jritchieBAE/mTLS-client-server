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

func main() {

	server, err := mtls.NewTLSServer(":8443", "../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hello", helloHandler)

	log.Fatal(server.Listen())
}
