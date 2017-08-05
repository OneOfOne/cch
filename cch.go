package cch

import (
	"errors"
)

var (
	ErrClosed = errors.New("channel is closed")
)

type Chan struct {
	w chan interface{}
	c chan struct{}
}

func Make(cap int) *Chan {
	ch := &Chan{
		w: make(chan interface{}, cap),
		c: make(chan struct{}),
	}
	return ch
}

func (ch *Chan) Closed() bool {
	select {
	case <-ch.c:
		return true
	default:
		return false
	}
}

func (ch *Chan) Send(val interface{}) (v bool) {
	select {
	case <-ch.c:
		break
	case ch.w <- val:
		v = true
	}
	return
}

func (ch *Chan) TrySend(val interface{}) (v bool) {
	select {
	case <-ch.c:
		break
	case ch.w <- val:
		v = true
	default:
	}
	return
}

func (ch *Chan) Recv() (v interface{}, ok bool) {
	select {
	case <-ch.c:
		break
	case v, ok = <-ch.w:
	}
	return
}

func (ch *Chan) TryRecv() (v interface{}, ok bool) {
	select {
	case <-ch.c:
		break
	case v, ok = <-ch.w:
	default:
	}
	return
}

func (ch *Chan) Chan(cap int) <-chan interface{} {
	out := make(chan interface{}, cap)
	go func() {
	L:
		for {
			select {
			case <-ch.c:
				break L
			case v, ok := <-ch.w:
				if !ok {
					break L
				}
				out <- v
			}
		}
		close(out)
	}()
	return out
}

func (ch *Chan) Close() error {
	select {
	case <-ch.c:
		return ErrClosed
	default:
		close(ch.c)
	}

	for {

		select {
		case <-ch.w:
		default:
			return nil
		}
	}

	return nil
}
