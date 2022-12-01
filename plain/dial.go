package plain

import (
	"net"
	"syscall"
	"time"
)

func DailTimeout(network, addr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, address string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, addr)
}
