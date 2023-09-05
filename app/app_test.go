package app

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestIPGen(t *testing.T) {
	workers := runtime.NumCPU()
	randIntCh := make(chan uint32, workers)
	defer close(randIntCh)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case r := <-randIntCh:
				b1 := uint8(r & 0xff)
				b2 := uint8((r >> 8) & 0xff)
				b3 := uint8((r >> 16) & 0xff)
				b4 := uint8(r >> 24)
				ip := net.IPv4(b4, b3, b2, b1)
				fmt.Println(ip.String())
			case <-ctx.Done():
				if len(randIntCh) == 0 {
					return
				}
			}
		}(i)
	}

	s := rand.NewSource(time.Now().UnixMicro())

	for i := 0; i < 20; i++ {
		genRandInt(randIntCh, &s)
	}
	cancel()
	wg.Wait()
}
