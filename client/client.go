package main

import (
	"fmt"
	"io/ioutil"
	"log"

	mtls "../pkg"
)

func main() {

	client, err := mtls.NewTLSClient("../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	}

	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", body)
}
