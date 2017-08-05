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

// Send send val to the channel, returns false if the channel is closed or block is false and the channel is full.
func (ch *Chan) Send(val interface{}, block bool) (ok bool) {
	ch.m.RLock()
	if ch.q != closedChan {
		if block {
			select {
			case <-ch.done:
			case ch.q <- val:
				ok = true
			}
		} else {
			select {
			case <-ch.done:
			case ch.q <- val:
				ok = true
			default:
			}
		}
	}
	ch.m.RUnlock()
	return
}

// Recv reads from the channel and blocks if block is set, until a value is available or the channel is closed.
func (ch *Chan) Recv(block bool) (v interface{}, ok bool) {
	if block {
		v, ok = <-ch.ch()
	} else {
		select {
		case v, ok = <-ch.ch():
		default:
		}
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
// Note that the call will block until the channel is fully drained.
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
	drain(w)
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
