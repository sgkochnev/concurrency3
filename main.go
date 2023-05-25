package main

import (
	"fmt"
)

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func (r *ringBuffer) Run() {
	defer close(r.outCh)

	if cap(r.outCh) == 0 {
		for _ = range r.inCh {
		}
	}

	for val := range r.inCh {
		select {
		case r.outCh <- val:
		default:
			<-r.outCh
			r.outCh <- val
		}
	}
}

func main() {
	inCh := make(chan int)
	outCh := make(chan int, 0)
	rb := NewRingBuffer(inCh, outCh)
	go rb.Run()

	max := 100
	for i := 0; i < max; i++ {
		inCh <- i
	}
	close(inCh)

	resSlice := make([]int, 0)
	for res := range outCh {
		resSlice = append(resSlice, res)
	}
	fmt.Println(resSlice)
}
