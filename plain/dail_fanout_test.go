package plain

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancleFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	t.Logf("listen on %s", listener.Addr().String())
	go func() {
		// Only accepting a single connection
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	dialer := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup) {
		defer wg.Done()
		t.Logf("dialer %d", id)
		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			return
		}
		c.Close()

		select {
		case <-ctx.Done():
		case response <- id:
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dialer(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res

	cancel()
	wg.Wait()

	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %v", ctx.Err())
	}
	t.Logf("dialer %d retrieved the resource", response)
}
