package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const rootCert = "/certwork/ca.crt"

func loadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
func main() {
	roots := x509.NewCertPool()
	home := os.Getenv("HOME")
	cert, err := loadFile(home + rootCert)
	if err != nil {
		log.Fatal(err)
	}
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		log.Fatal("failed to parse root certificate")
	}
	config := &tls.Config{RootCAs: roots, ServerName: "*.nas.local"}

	connp, err := net.Dial("tcp", "test.nas.local:8089")
	if err != nil {
		log.Fatal(err)
	}

	conn := tls.Client(connp, config)
	io.WriteString(conn, "Hello secure Server")
	conn.Close()
}
