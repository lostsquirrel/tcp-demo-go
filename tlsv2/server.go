package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const certPath = "/certwork/nas/server.crt"
const keyPath = "/certwork/nas/server.key"

func loadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func main() {
	home := os.Getenv("HOME")
	cert, err := loadFile(home + certPath)
	if err != nil {
		log.Fatal(err)
	}
	key, err := loadFile(home + keyPath)
	cer, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	lp, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal(err)
	}
	defer lp.Close()

	l := tls.NewListener(lp, config)
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			fmt.Println()
			c.Close()
		}(conn)
	}
}
