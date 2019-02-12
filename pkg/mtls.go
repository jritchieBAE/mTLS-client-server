package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

type tlsServer struct {
	server   *http.Server
	certPath string
	keyPath  string
}

func (t *tlsServer) Listen() error {
	return t.server.ListenAndServeTLS(t.certPath, t.keyPath)
}

func (t *tlsServer) ListenNoTLS() error {
	return t.server.ListenAndServe()
}

func NewTLSServer(listenAddress string, certPath, keyPath string) (*tlsServer, error) {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		GetCertificate: func(h *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(certPath, keyPath)
			return &cert, err
		},
	}

	server := &http.Server{
		Addr:      listenAddress,
		TLSConfig: tlsConfig,
	}

	tls := &tlsServer{
		server:   server,
		certPath: certPath,
		keyPath:  keyPath,
	}

	return tls, nil
}

type tlsClient struct {
	client   *http.Client
	certPath string
	keyPath  string
}

func (t *tlsClient) Get(path string) (resp *http.Response, err error) {
	return t.client.Get(path)
}

func NewTLSClient(certPath, keyPath string) (*tlsClient, error) {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
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

	t := &tlsClient{
		client:   client,
		certPath: certPath,
		keyPath:  keyPath,
	}
	return t, nil
}
