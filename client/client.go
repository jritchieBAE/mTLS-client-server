package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"../mtls"
)

var (
	certPath = "cert.pem"
	keyPath  = "key.pem"
)

type ClientInterface interface {
	Get(string) (*http.Response, error)
}

func main() {
	urls := [...]string{
		// "http://localhost:8443/hello",
		// "https://localhost:8443/hello",
		//"https://localhost:8443/hello",
		"https://localhost:9100/metrics",
	}

	functions := [...]func(string){
		unsecuredClient,
		tlsClient,
		mtlsClient,
	}
	for _, f := range functions {
		for _, url := range urls {
			fmt.Printf("\nTesting connection with %s over %s\n", fName(f), url)
			f(url)
		}
	}
}

func unsecuredClient(url string) {
	client := mtls.NewUnsecureClient()
	r, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	} else {

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", body)
	}
}

func tlsClient(url string) {

	client, err := mtls.NewTlsClient(certPath)
	if err != nil {
		log.Fatal(err)
	}
	r, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	} else {

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s\n", body)
	}
}

func mtlsClient(url string) {

	client, err := mtls.NewMtlsClient(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	r, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
	} else {

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s\n", body)
	}
}

func fName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
