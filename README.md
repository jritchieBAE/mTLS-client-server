# mTLS-client-server

Example implementation of mTLS in GoLang

Based on this tutorial: https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go

The mtls package provides the following APIs:

To create a server:
- NewUnsecureServer()
- NewTlsServer(certPath, keyPath string)
- NewMtlsServer(certPath, keyPath string)

The server object returned by these methods can be started with a call to server.Listen(address string)


To create a client:
- NewUnsecureClient()
- NewTlsClient(certPath string)
- NewMtlsClient(cerPath, keyPath string)

The client object returned by these methods can read from a server with a call to client.Get(url)