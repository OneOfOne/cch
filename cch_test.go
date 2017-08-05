package cch

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
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
		var (
			last interface{}
		)
		for v := range ch.Chan() {
			if v := v.(int64); v == 1e4 {
				ch.Close()
			} else {
				last = v
			}
		}
		_ = last
		// t.Logf("done: %v", last)

	}()
	wg.Wait()
	if _, ok := ch.Recv(); ok {
		t.Fatal("unexpected Recv")
	}

}

func BenchmarkTryLeak(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var (
			ch = Make(runtime.NumCPU())
			i  int
		)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for range ch.Chan() {

				}
			}()
		}
		for pb.Next() {
			if !ch.Send(i) {
				b.Fatalf("wtf")
			}
			i++
		}
		ch.Close()
		b.Log(i)
	})
}
