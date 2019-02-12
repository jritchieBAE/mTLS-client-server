package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

type tlsServer struct {
	server   *http.Server
	port     string
	certPath string
	keyPath  string
}

func (t *tlsServer) Listen() error {
	return t.server.ListenAndServeTLS(t.certPath, t.keyPath)
}

func (t *tlsServer) SetPort(port string) {
	t.port = ":" + port
	t.server.Addr = t.port
}

func NewTLSServer(certPath, keyPath string) (*tlsServer, error) {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
	tls := &tlsServer{
		server:   server,
		port:     server.Addr,
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
