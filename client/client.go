package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	mtls "../pkg"
)

const (
	htp  = "http://"
	htps = "https://"
	url  = "localhost:8443/hello"
)

type ClientInterface interface {
	Get(string) (*http.Response, error)
}

func main() {

	fmt.Println("\nTesting connection with unsecured over http")
	r, err := http.Get(htp + url)
	if err != nil {
		log.Println(err)
	} else {
		output(r)
	}

	fmt.Println("\nTesting connection with unsecured over https")
	r, err = http.Get(htps + url)
	if err != nil {
		log.Println(err)
	} else {
		output(r)
	}

	fmt.Println("\nTesting connection with TLS over http")
	client, err := mtls.NewTLSClient("../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	} else {
		readFrom(client, htp+url)
	}

	fmt.Println("\nTesting connection with TLS over https")
	client, err = mtls.NewTLSClient("../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	} else {
		readFrom(client, htps+url)
	}

	fmt.Println("\nTesting connection with mTLS over http")
	client, err = mtls.NewMTLSClient("../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	} else {
		readFrom(client, htp+url)
	}

	fmt.Println("\nTesting connection with mTLS over https")
	client, err = mtls.NewMTLSClient("../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	} else {
		readFrom(client, htps+url)
	}
}

func readFrom(client ClientInterface, url string) {
	r, err := client.Get(url)
	if err != nil {
		log.Println(err)
	} else {
		output(r)
	}
}
func output(r *http.Response) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", body)
}
