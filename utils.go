package cch

import (
	"runtime"
)

// SelectSend will send val on any of the channels, in order.
// it returns false if block is false and all channels are full or all channels are closed.
func SelectSend(val interface{}, block bool, chans ...*Chan) (ok bool) {
	if len(chans) == 0 {
		return
	}
	for {
		numClosed := 0
		for _, ch := range chans {
			if ch.Send(val, false) {
				return true
			}
			if ch.Closed() {
				numClosed++
			}
		}

		if numClosed == len(chans) || !block {
			return false
		}

		runtime.Gosched()
	}

}

// SelectRecv will read val from any of the input channels, in order.
func SelectRecv(block bool, chans ...*Chan) (val interface{}, ok bool) {
	if len(chans) == 0 {
		return
	}
	for {
		numClosed := 0
		for _, ch := range chans {
			if val, ok = ch.Recv(false); ok {
				return
			}

			if ch.Closed() {
				numClosed++
			}
		}

		if numClosed == len(chans) || !block {
			return
		}

		runtime.Gosched()
	}

}

func CloseAll(chans ...*Chan) {
	for _, ch := range chans {
		ch.Close()
	}
}
