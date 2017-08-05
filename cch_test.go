package cch

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestChan(t *testing.T) {
	var (
		ch = NewBuffered(2)
		wg sync.WaitGroup
		i  int64
	)
	wg.Add(3)
	go func() {
		defer wg.Done()
		for {
			if !ch.Send(atomic.AddInt64(&i, 1), true) {
				//t.Logf("Send1: chan closed at %d", i)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			if !ch.Send(atomic.AddInt64(&i, 1), true) {
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
		t.Logf("done: %v", last)

	}()
	wg.Wait()
	if _, ok := ch.Recv(true); ok {
		t.Fatal("unexpected Recv")
	}

}

func TestSelect(t *testing.T) {
	const N = 1e3
	chans := make([]*Chan, N)
	for i := range chans {
		chans[i] = NewBuffered(1)
	}

	for i := range chans {
		if !SelectSend(i, false, chans...) {
			t.Fatalf("couldn't send %v", i)
		}
	}

	for i := range chans {
		if _, ok := SelectRecv(false, chans...); !ok {
			t.Fatalf("couldn't recv %v", i)
		}
	}
}

func BenchmarkTryLeak(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var (
			ch = NewBuffered(runtime.NumCPU())
			i  int
		)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for range ch.Chan() {

				}
			}()
		}
		for pb.Next() {
			if !ch.Send(i, true) {
				b.Fatalf("wtf")
			}
			i++
		}
		ch.Close()
		b.Log(i)
	})
}
