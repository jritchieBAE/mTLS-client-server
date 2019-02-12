package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

type TLSServer struct {
	_server   *http.Server
	_port     string
	_certPath string
	_keyPath  string
}

func (t *TLSServer) Listen() {
	log.Fatal(t._server.ListenAndServeTLS(t._certPath, t._keyPath))
}

func (t *TLSServer) SetPort(port string) {
	t._port = ":" + port
}

func NewTLSServer(certPath, keyPath string) *TLSServer {
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
	tls := &TLSServer{
		_server:   server,
		_port:     server.Addr,
		_certPath: certPath,
		_keyPath:  keyPath,
	}

	return tls
}
