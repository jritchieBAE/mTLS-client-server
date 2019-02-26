package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"

	mtls "../mtlsClient"
)

var (
	certPath   = "../cert.crt"
	keyPath    = "../cert.key"
	caCertPath = "../root.crt"
)

// client.go tests a URL to see if it can be connected to by unsecured HTTP,
// TLS, or mTLS, using the provided certificates. An alternative URL to
// test can be provided as a command line parameter.
func main() {
	url := "http://localhost:8080"

	if len(os.Args) > 1 {
		args := os.Args[1:]
		re := regexp.MustCompile(`http.?:\/\/.*`)
		if re.Match([]byte(args[0])) {
			url = args[0]
		} else {
			log.Fatal("URL not recognised: " + args[0])
		}
	}

	functions := [...]func(string){
		unsecuredClient,
		tlsClient,
		mtlsClient,
	}
	for _, f := range functions {
		fmt.Printf("\nTesting connection with %s over %s\n", fName(f), url)
		f(url)
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

	client, err := mtls.NewTlsClient(caCertPath)
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

	client, err := mtls.NewMtlsClient(certPath, keyPath, caCertPath)
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
