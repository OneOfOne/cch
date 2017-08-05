package cch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	var (
		ch = Make(2)
		wg sync.WaitGroup
		i  int64
	)
	wg.Add(3)
	go func() {
		defer wg.Done()
		for {
			if !ch.Send(atomic.AddInt64(&i, 1)) {
				//t.Logf("Send1: chan closed at %d", i)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			if !ch.Send(atomic.AddInt64(&i, 1)) {
				//t.Logf("Send2: chan closed at %d", i)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		var (
			last interface{}
		)
		for v := range ch.Chan(0) {
			if v := v.(int64); v == 1e6 {
				ch.Close()
			} else {
				last = v
			}
		}
		t.Logf("done: %v", last)

	}()
	wg.Wait()
	t.Log(ch.Recv())

}

func BenchmarkLeak(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var (
			ch   = Make(1)
			last interface{}
			i    int
		)
		go func() {
			for v := range ch.Chan(0) {
				last = v
			}
			b.Logf("last: %v", last)
		}()
		for pb.Next() {
			if !ch.Send(i) {
				b.Fatalf("wtf")
			}
			i++
		}
		ch.Close()
	})
}
