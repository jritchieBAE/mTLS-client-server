package mtlsServer

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

func NewMtlsServer(certPath, keyPath, caCertPath string) (*TlsServer, error) {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
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

	tlsConfig := &tls.Config{
		ClientCAs:    clientCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
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
