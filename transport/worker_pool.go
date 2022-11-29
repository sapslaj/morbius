package transport

import (
	"sync"
)

type WorkerPool[V any] struct {
	WorkerCount    int
	MessageHandler func(V)
	MessageChannel chan V
	wg             sync.WaitGroup
}

func NewWorkerPool[V any](workerCount int, messageBuffer int, handler func(V)) *WorkerPool[V] {
	messageChannel := make(chan V, messageBuffer)
	return &WorkerPool[V]{
		WorkerCount:    workerCount,
		MessageHandler: handler,
		MessageChannel: messageChannel,
	}
}

func (wp *WorkerPool[V]) Start() {
	for i := 0; i < wp.WorkerCount; i++ {
		wp.wg.Add(1)
		go func(i int) {
			defer wp.wg.Done()
			for msg := range wp.MessageChannel {
				wp.MessageHandler(msg)
			}
		}(i)
	}
}

func (wp *WorkerPool[V]) Push(message V) {
	wp.MessageChannel <- message
}
