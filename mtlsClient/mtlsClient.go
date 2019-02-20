package mtlsClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

func NewMtlsClient(certPath, keyPath, caCertPath string) (*http.Client, error) {

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

	return client, nil
}

func NewTlsClient(CAcertPath string) (*http.Client, error) {
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

	return client, nil
}

func NewUnsecureClient() *http.Client {

	return http.DefaultClient
}
