package cch

import (
	"errors"
	"sync"
)

var (
	// ErrClosed is returned when Chan.Closed() gets called multiple times.
	ErrClosed = errors.New("channel is closed")

	closedChan = make(chan interface{})
)

func init() {
	close(closedChan) // inspired by context
}

// Chan is a closable channel.
type Chan struct {
	q    chan interface{}
	done chan struct{}
	m    sync.RWMutex
}

// New returns a new unbuffed channel, same as NewBuffered(0).
func New() *Chan { return NewBuffered(0) }

// NewBuffered returns a new closable channel with the specific cap, use 0 for an unbuffered channel.
func NewBuffered(cap int) *Chan {
	ch := &Chan{
		q:    make(chan interface{}, cap),
		done: make(chan struct{}),
	}
	return ch
}

// Closed returns whether the channel is closed or not.
func (ch *Chan) Closed() bool {
	select {
	case <-ch.done:
		return true
	default:
		return false
	}
}

// Send send val to the channel, returns false if the channel is closed.
func (ch *Chan) Send(val interface{}) (ok bool) {
	ch.m.RLock()
	if ch.q != closedChan {
		select {
		case <-ch.done:
		case ch.q <- val:
			ok = true
		}
	}
	ch.m.RUnlock()
	return
}

// TrySend tries to send val to the channel, returns false if the channel is closed or the channel is full.
func (ch *Chan) TrySend(val interface{}) (ok bool) {
	ch.m.RLock()
	if ch.q != closedChan {
		select {
		case ch.q <- val:
			ok = true
		default:
		}
	}
	ch.m.RUnlock()
	return
}

// Recv reads from the channel and blocks until a value is available or the channel is closed.
func (ch *Chan) Recv() (v interface{}, ok bool) {
	v, ok = <-ch.ch()
	return
}

// Recv tries to read from the channel and returns early if the channel is empty or closed.
func (ch *Chan) TryRecv() (v interface{}, ok bool) {
	select {
	case v, ok = <-ch.ch():
	default:
	}
	return
}

// Chan returns a read-only version of the underlying channel.
func (ch *Chan) Chan() <-chan interface{} {
	return ch.ch()
}

// Len returns the length of the channel.
func (ch *Chan) Len() int { return len(ch.ch()) }

// Len returns the capacity of the channel.
func (ch *Chan) Cap() int { return cap(ch.ch()) }

// Close closes the channel and drains it safely.
// Close returns ErrClosed if called multiple times.
func (ch *Chan) Close() error {
	select {
	case <-ch.done:
		return ErrClosed
	default:
		close(ch.done)
	}

	ch.m.Lock()
	w := ch.q
	ch.q = closedChan
	ch.m.Unlock()
	go drain(w)
	return nil
}

func (ch *Chan) ch() chan interface{} {
	ch.m.RLock()
	w := ch.q
	ch.m.RUnlock()
	return w
}

func drain(w chan interface{}) {
	for {
		select {
		case <-w:
		default:
			close(w)
			return
		}
	}
}
