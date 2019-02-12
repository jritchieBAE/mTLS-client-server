package main

import (
	"io"
	"net/http"

	mtls "../pkg"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!\n")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	server := mtls.NewTLSServer("../cert.pem", "../key.pem")

	server.Listen()
}
