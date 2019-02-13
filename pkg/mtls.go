package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TlsServer struct {
	listen (func(string) error)
}

func (t *TlsServer) Listen(address string) error {
	if t.listen != nil {
		return t.listen(address)
	} else {
		return errors.New("Server not instantiated")
	}
}

func NewMtlsServer(cert, key string) *TlsServer {

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

	l := func(addr string) error {
		server := &http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
		}
		return server.ListenAndServeTLS(cert, key)
	}

	return &TlsServer{listen: l}
}

func NewTlsServer(cert, key string) *TlsServer {

	l := func(addr string) error {
		return http.ListenAndServeTLS(addr, cert, key, nil)
	}

	return &TlsServer{listen: l}
}

func NewUnsecureServer() *TlsServer {

	l := func(addr string) error {
		return http.ListenAndServe(addr, nil)
	}

	return &TlsServer{listen: l}
}

type TlsClient struct {
	get func(string) (*http.Response, error)
}

func (t *TlsClient) Get(url string) (*http.Response, error) {
	return t.get(url)
}

func NewMtlsClient(certPath, keyPath string) *TlsClient {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Println(err)
	}

	caCert, err := ioutil.ReadFile(certPath)
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

	g := func(url string) (*http.Response, error) {
		return client.Get(url)
	}

	return &TlsClient{get: g}
}

func NewTlsClient(cert string) *TlsClient {
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

	g := func(url string) (*http.Response, error) {
		return client.Get(url)
	}

	return &TlsClient{get: g}
}

func NewUnsecureClient() *TlsClient {
	g := func(url string) (*http.Response, error) {
		return http.Get(url)
	}

	return &TlsClient{get: g}
}
