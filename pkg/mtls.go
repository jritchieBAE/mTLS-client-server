package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

type tlsServer struct {
	_server   *http.Server
	_port     string
	_certPath string
	_keyPath  string
}

func (t *tlsServer) Listen() {
	log.Fatal(t._server.ListenAndServeTLS(t._certPath, t._keyPath))
}

func (t *tlsServer) SetPort(port string) {
	t._port = ":" + port
}

func NewTLSServer(certPath, keyPath string) *tlsServer {
	caCert, err := ioutil.ReadFile(certPath)
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

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
	tls := &tlsServer{
		_server:   server,
		_port:     server.Addr,
		_certPath: certPath,
		_keyPath:  keyPath,
	}

	return tls
}

type tlsClient struct {
	_client   *http.Client
	_certPath string
	_keyPath  string
}

func (t *tlsClient) Get(path string) (resp *http.Response, err error) {
	return t._client.Get(path)
}

func NewTLSClient(certPath, keyPath string) *tlsClient {

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatal(err)
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
		_client:   client,
		_certPath: certPath,
		_keyPath:  keyPath,
	}
	return t
}
