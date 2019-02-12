package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

type tlsServer struct {
	server    *http.Server
	tlsConfig *tls.Config
	certPath  string
	keyPath   string
}

func (t *tlsServer) ListenMTLS(listenAddress string) error {
	if len(t.certPath) > 0 {
		t.server = &http.Server{
			Addr:      listenAddress,
			TLSConfig: t.tlsConfig,
		}
		return t.server.ListenAndServeTLS(t.certPath, t.keyPath)
	} else {
		return errors.New("Cannot start mTLS server as no certificates have been defined")
	}
}

func (t *tlsServer) ListenTLS(listenAddress string) error {
	if len(t.certPath) > 0 {
		return http.ListenAndServeTLS(listenAddress, t.certPath, t.keyPath, nil)
	} else {
		return errors.New("Cannot start TLS server as no certificates have been defined")
	}
}

func (t *tlsServer) ListenNoTLS(listenAddress string) error {
	t.server = &http.Server{
		Addr: listenAddress,
	}
	return t.server.ListenAndServe()
}

func (t *tlsServer) WithCertificates(certPath, keyPath string) (*tlsServer, error) {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t.certPath, t.keyPath = certPath, keyPath

	t.tlsConfig = &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		GetCertificate: func(h *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(certPath, keyPath)
			return &cert, err
		},
	}

	return t, nil
}

func NewTLSServer() *tlsServer {
	tls := &tlsServer{}
	return tls
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
