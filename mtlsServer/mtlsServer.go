package mtlsServer

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

func NewMtlsServer(addr, certPath, keyPath, caCertPath string) (*http.Server, error) {

	server, err := NewTlsServer(addr, certPath, keyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	clientCAs := x509.NewCertPool()
	if ok := clientCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("Failed to parse CA Certificate")
	}

	server.TLSConfig.ClientCAs = clientCAs
	server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert

	return server, nil
}

func NewTlsServer(addr, certPath, keyPath string) (*http.Server, error) {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	return server, nil
}

func NewServer(addr string) *http.Server {

	return &http.Server{
		Addr: addr,
	}
}
