# mTLS-client-server

Example implementation of mTLS in GoLang

Originally based on this tutorial: https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go

The mtlsServer package provides the following APIs:

To create a server:
- NewUnsecureServer()
- NewTlsServer(certPath, keyPath string)
- NewMtlsServer(certPath, keyPath string)

The object (*TlsServer) returned by these methods can be started with a call to server.Listen(address string). This provides a standard interface rather than specifying whether the server is to Listen, 


To create a client:
- NewUnsecureClient()
- NewTlsClient(certPath string)
- NewMtlsClient(cerPath, keyPath string)

The object returned by these methods is a standard http.Client object and can be used to contact a server as per the standard API