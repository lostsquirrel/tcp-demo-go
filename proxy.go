package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Proxy struct {
	address          string
	target           *net.TCPAddr
	terminationDelay time.Duration
}

func NewProxy(address string, terminationDelay time.Duration) (*Proxy, error) {
	_, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		address:          address,
		terminationDelay: terminationDelay,
	}, nil

}

func (p *Proxy) ServeTCP(conn net.Conn) {
	log.Printf("Handling connection from %s", conn.RemoteAddr().String())

	connBackend, err := p.dialBackend()
	if err != nil {
		log.Printf("Error while connecting to backend: %v", err)
		return
	}

	// maybe not needed, but just in case
	defer connBackend.Close()
	errChan := make(chan error)
	go p.connCopy(conn, connBackend, errChan)
	go p.connCopy(connBackend, conn, errChan)

	err = <-errChan
	if err != nil {
		log.Printf("Error during connection: %v", err)
	}

	<-errChan
}

func (p Proxy) dialBackend() (*tls.Conn, error) {

	roots, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
		roots = x509.NewCertPool()
	}
	log.Println("load system certificates")
	config := &tls.Config{RootCAs: roots}

	conn, err := tls.Dial("tcp", p.address, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (p Proxy) connCopy(dst, src net.Conn, errCh chan error) {
	_, err := io.Copy(dst, src)
	log.Println(err)
	errCh <- err

	if p.terminationDelay >= 0 {
		err := dst.SetReadDeadline(time.Now().Add(p.terminationDelay))
		if err != nil {
			log.Printf("Error while setting deadline: %v", err)
		}
	}
}

func main() {
	port := flag.Int("port", 9000, "Server Port")
	backend := flag.String("backend", "localhost:8443", "backend address")
	flag.Parse()
	log.Println("get backend " + *backend)
	proxy := Proxy{
		address: *backend,
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			proxy.ServeTCP(conn)
			c.Close()
		}(conn)
	}
}
