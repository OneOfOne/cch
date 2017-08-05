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
	m sync.RWMutex
}

func Make(cap int) *Chan {
	ch := &Chan{
		w: make(chan interface{}, cap),
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
		ch.w <- val
		ok = true

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
	v, ok = <-ch.w
	return
}

func (ch *Chan) TryRecv() (v interface{}, ok bool) {
	select {
	case v, ok = <-ch.w:
	default:
	}
	return
}

func (ch *Chan) Chan(cap int) <-chan interface{} {
	return ch.ch()
}

func (ch *Chan) Close() error {
	ch.m.Lock()
	if ch.w == nil {
		ch.m.Unlock()
		return ErrClosed
	}

L:
	for {
		select {
		case <-ch.w:
		default:
			close(ch.w)
			break L
		}
	}
	ch.w = closedChan
	ch.m.Unlock()
	return nil
}
