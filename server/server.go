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

	http.HandleFunc("/hello", helloHandler)

	// server, err := (mtls.NewTLSServer().WithCertificates("../cert.pem", "../key.pem"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Fatal(server.Listen(":8443"))

	server := mtls.NewTLSServer()
	log.Fatal(server.Listen(":8080"))
}
