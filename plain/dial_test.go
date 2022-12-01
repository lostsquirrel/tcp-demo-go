package plain

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("listen on %s", listener.Addr().String())
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}
			t.Logf("accept from %s", conn.RemoteAddr().String())
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}

			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	time.Sleep(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	<-done
	listener.Close()
	<-done
}
