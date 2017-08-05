package cch

import (
	"errors"
	"sync"
)

var (
	ErrClosed  = errors.New("channel is closed")
	closedChan = make(chan interface{})
)

func init() {
	close(closedChan) // inspired by context
}

type Chan struct {
	w chan interface{}
	c chan struct{}
	m sync.RWMutex
}

func Make(cap int) *Chan {
	ch := &Chan{
		w: make(chan interface{}, cap),
		c: make(chan struct{}),
	}
	return ch
}

func (ch *Chan) ch() chan interface{} {
	ch.m.RLock()
	w := ch.w
	ch.m.RUnlock()
	return w
}

func (ch *Chan) Closed() bool {
	return ch.ch() == closedChan
}

func (ch *Chan) Send(val interface{}) (ok bool) {
	ch.m.RLock()
	if ch.w != closedChan {
		select {
		case <-ch.c:
		case ch.w <- val:
			ok = true
		}
	}
	ch.m.RUnlock()
	return
}

func (ch *Chan) TrySend(val interface{}) (ok bool) {
	if w := ch.ch(); w != closedChan {
		select {
		case ch.w <- val:
			ok = true
		default:
		}
	}
	return
}

func (ch *Chan) Recv() (v interface{}, ok bool) {
	v, ok = <-ch.ch()
	return
}

func (ch *Chan) TryRecv() (v interface{}, ok bool) {
	select {
	case v, ok = <-ch.ch():
	default:
	}
	return
}

func (ch *Chan) Chan() <-chan interface{} {
	return ch.ch()
}

func (ch *Chan) Close() error {
	select {
	case <-ch.c:
		return ErrClosed
	default:
		close(ch.c)
	}
	ch.m.Lock()
	w := ch.w
	ch.w = closedChan
	ch.m.Unlock()
	drain(w)
	return nil
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
