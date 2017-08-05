package cch

import (
	"runtime"
	"sync"
)

type MultiChan struct {
	chans map[*Chan]struct{}
	m     sync.RWMutex
}

func Select(chans ...*Chan) *MultiChan {
	out := make(map[*Chan]struct{}, len(chans))

	for _, ch := range chans {
		out[ch] = struct{}{}
	}

	return &MultiChan{chans: out}
}

func NewMulti(numChans, bufferPerChan int) *MultiChan {
	out := make(map[*Chan]struct{}, numChans)

	for i := 0; i < numChans; i++ {
		out[NewBuffered(bufferPerChan)] = struct{}{}
	}

	return &MultiChan{chans: out}
}

func (mc *MultiChan) delete(old *Chan) {
	new := make(map[*Chan]struct{})

	mc.m.Lock()

	for ch := range mc.chans {
		new[ch] = struct{}{}
	}

	delete(new, old)

	mc.chans = new

	mc.m.Unlock()
}

// Send will send to any of the underlying channels in a random order.
// if the channel is closed externally, it will be removed from MultiChan.
func (mc *MultiChan) Send(val interface{}, block bool) (ok bool) {
	for {
		mc.m.RLock()
		chans := mc.chans
		mc.m.RUnlock()

		if len(chans) == 0 {
			return
		}

		for ch := range chans {
			if ch.Send(val, false) {
				return true
			}

			if ch.Closed() {
				mc.delete(ch)
			}
		}

		if !block {
			return
		}

		runtime.Gosched()
	}
}

func (mc *MultiChan) Recv(block bool) (val interface{}, ok bool) {
	for {
		mc.m.RLock()
		chans := mc.chans
		mc.m.RUnlock()

		if len(chans) == 0 {
			return
		}

		for ch := range chans {
			if val, ok = ch.Recv(false); ok {
				return
			}

			if ch.Closed() {
				mc.delete(ch)
			}
		}

		if !block {
			return
		}

		runtime.Gosched()
	}
}

// Chan returns a read-only version of the underlying channel(s).
// this function will leak a goroutine until mc.Close() is called
func (mc *MultiChan) Chan() <-chan interface{} {
	mc.m.RLock()
	ln := len(mc.chans)
	mc.m.RUnlock()
	out := make(chan interface{}, ln)
	go func() {
		for {
			v, ok := mc.Recv(true)
			if !ok {
				break
			}
			out <- v
		}
		close(out)
	}()
	return out
}

func (mc *MultiChan) Close() error {
	mc.m.Lock()
	chans := mc.chans
	mc.chans = nil
	mc.m.Unlock()

	for ch := range chans {
		ch.Close()
	}

	return nil
}
