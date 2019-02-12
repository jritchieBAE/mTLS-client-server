package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

const ()

type ClientInterface interface {
	Get(string) (*http.Response, error)
}

func main() {
	urls := [...]string{
		"http://localhost:8443/hello",
		"https://localhost:8443/hello",
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
	r, err := http.Get(url)
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

	caCert, err := ioutil.ReadFile("../cert.pem")
	if err != nil {
		fmt.Println(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
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

	cert, err := tls.LoadX509KeyPair("../cert.pem", "../key.pem")
	if err != nil {
		fmt.Println(err)
	}

	caCert, err := ioutil.ReadFile("../cert.pem")
	if err != nil {
		fmt.Println(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
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
