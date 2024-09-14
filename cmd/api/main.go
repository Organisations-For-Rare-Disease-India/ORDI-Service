package main

import (
	"ORDI/internal/server"
	"fmt"
	"os"
)

func main() {

	server := server.NewServer()
	certFile := os.Getenv("SSL_CERT_PATH")
	keyFile := os.Getenv("SSL_KEY_PATH")
	isLocal := os.Getenv("IS_LOCAL") == "true"
	var err error
	if isLocal {
		err = server.ListenAndServe()
	} else {
		server.Addr = ":443"
		err = server.ListenAndServeTLS(certFile, keyFile)
	}
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
