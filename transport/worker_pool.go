package transport

import (
	"errors"
	"fmt"
	"sync"
)

type WorkerPool[V any] struct {
	WorkerCount    int
	MessageHandler func(V)
	MessageChannel chan V
	wg             sync.WaitGroup
	messageBuffer  int
}

func NewWorkerPool[V any](workerCount int, messageBuffer int, handler func(V)) *WorkerPool[V] {
	wp := &WorkerPool[V]{
		WorkerCount:    workerCount,
		MessageHandler: handler,
		messageBuffer:  messageBuffer,
	}
	wp.openMessageChannel()
	return wp
}

func (wp *WorkerPool[V]) openMessageChannel() {
	if wp.MessageChannel == nil {
		wp.MessageChannel = make(chan V, wp.messageBuffer)
	}
}

func (wp *WorkerPool[V]) Start() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()
	wp.openMessageChannel()
	for i := 0; i < wp.WorkerCount; i++ {
		wp.wg.Add(1)
		go func(i int) {
			defer wp.wg.Done()
			for msg := range wp.MessageChannel {
				wp.MessageHandler(msg)
			}
		}(i)
	}
	return
}

func (wp *WorkerPool[V]) Stop() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()
	close(wp.MessageChannel)
	wp.wg.Wait()
	return
}

func (wp *WorkerPool[V]) Push(message V) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()
	wp.MessageChannel <- message
	return
}
