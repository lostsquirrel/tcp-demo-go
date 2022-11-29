package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
)

func main() {
	roots, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
		roots = x509.NewCertPool()
	}
	config := &tls.Config{RootCAs: roots}

	conn, err := tls.Dial("tcp", "home.shangao.tech:6690", config)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(conn, "Hello simple secure Server with system CA")
	conn.Close()
}
