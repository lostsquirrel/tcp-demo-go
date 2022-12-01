package start

import (
	"context"
	"fmt"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration
	select {
	case <-ctx.Done():
		return
	case interval = <-reset:
		fmt.Printf("current interval: %f\n", interval.Seconds())
	default:
	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)

	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case newInterval := <-reset:
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
			// fmt.Printf("received new interval: %d\n", interval.Round(100*time.Millisecond))
		case <-timer.C:
			// fmt.Println("send ping")
			if _, err := w.Write([]byte("ping")); err != nil {
				// track and act on consecutive timeouts here
				return
			}

		}
		_ = timer.Reset(interval)
	}
}
