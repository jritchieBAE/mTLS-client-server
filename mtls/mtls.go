package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
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

func NewMtlsServer(certPath, keyPath string) (*TlsServer, error) {

	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	_, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	l := func(addr string) error {
		server := &http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
		}
		return server.ListenAndServeTLS(certPath, keyPath)
	}

	return &TlsServer{listen: l}, nil
}

func NewTlsServer(certPath, keyPath string) (*TlsServer, error) {

	_, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	_, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	l := func(addr string) error {
		return http.ListenAndServeTLS(addr, certPath, keyPath, nil)
	}

	return &TlsServer{listen: l}, nil
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

func NewMtlsClient(certPath, keyPath string) (*TlsClient, error) {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	caCert, _ := ioutil.ReadFile(certPath)

	caCertPool, _ := x509.SystemCertPool()
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

	return &TlsClient{get: g}, nil
}

func NewTlsClient(certPath string) (*TlsClient, error) {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	caCertPool, _ := x509.SystemCertPool()
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

	return &TlsClient{get: g}, nil
}

func NewUnsecureClient() *TlsClient {
	g := func(url string) (*http.Response, error) {
		return http.Get(url)
	}

	return &TlsClient{get: g}
}
