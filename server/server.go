package main

import (
	"io"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!\n")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	log.Fatal(http.ListenAndServeTLS(":8443", "../cert.pem", "../key.pem", nil))
}
