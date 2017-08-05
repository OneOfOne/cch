package cch

// SelectSend will send val on any of the channels, in random order.
// it returns false if block is false and all channels are full or all channels are closed.
func SelectSend(val interface{}, block bool, chans ...*Chan) (ok bool) {
	return Select(chans...).Send(val, block)
}

// SelectRecv will read val from any of the input channels, in random order.
func SelectRecv(block bool, chans ...*Chan) (val interface{}, ok bool) {
	return Select(chans...).Recv(block)
}

func CloseAll(chans ...*Chan) {
	for _, ch := range chans {
		ch.Close()
	}
}
