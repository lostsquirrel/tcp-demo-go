package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func loadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

const rootCert = "/certwork/ca.crt"

func main() {
	roots := x509.NewCertPool()
	home := os.Getenv("HOME")
	cert, err := loadFile(home + rootCert)
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		log.Fatal("failed to parse root certificate")
	}
	config := &tls.Config{RootCAs: roots}

	conn, err := tls.Dial("tcp", "test.nas.local:6690", config)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(conn, "Hello simple secure Server")
	conn.Close()
}
