package mtlsClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

type TlsClient struct {
	get func(string) (*http.Response, error)
}

func (t *TlsClient) Get(url string) (*http.Response, error) {
	return t.get(url)
}

func NewMtlsClient(certPath, keyPath, caCertPath string) (*TlsClient, error) {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	clientCAs, _ := x509.SystemCertPool()
	if ok := clientCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("Failed to parse CA Certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:      clientCAs,
		Certificates: []tls.Certificate{cert},
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	g := func(url string) (*http.Response, error) {
		return client.Get(url)
	}

	return &TlsClient{get: g}, nil
}

func NewTlsClient(CAcertPath string) (*TlsClient, error) {
	caCert, err := ioutil.ReadFile(CAcertPath)
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
