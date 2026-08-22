package sse

import (
	"fmt"
	"sync"
)

type Broadcaster[T any] struct {
	mu            sync.Mutex
	output        chan T
	history       []T
	outBufferSize int
	maxSize       int
}

func (b *Broadcaster[T]) GetHistory() []T {
	return b.history
}

func NewBroadcaster[T any](historyMaxSize int, outBufferSize int) *Broadcaster[T] {
	return &Broadcaster[T]{
		output:        make(chan T, outBufferSize),
		history:       make([]T, 0),
		maxSize:       historyMaxSize,
		outBufferSize: outBufferSize,
	}
}

func (b *Broadcaster[T]) Publish(item T) {
	b.mu.Lock()
	b.history = append(b.history, item)
	if len(b.history) > b.maxSize {
		b.history = b.history[1:]
	}
	b.mu.Unlock()

	select {
	case b.output <- item:
	default:
	}
}

func (b *Broadcaster[T]) Subscribe() <-chan T {
	return b.output
}

func (b *Broadcaster[T]) RemoveAt(index int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if index < 0 || index >= len(b.history) {
		return fmt.Errorf("Out-of-bounds error: %d", index)
	}

	b.history = append(b.history[:index], b.history[index+1:]...)

	return nil
}

func (b *Broadcaster[T]) RemoveWhere(match func(T) bool) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	filtered := make([]T, 0)
	removed := 0

	for _, item := range b.history {
		if match(item) {
			removed++
			continue
		}
		filtered = append(filtered, item)
	}

	b.history = filtered

	return removed
}

func (b *Broadcaster[T]) RemoveOutput() {
	for {
		select {
		case <-b.Subscribe():
		default:
			return
		}
	}
}
