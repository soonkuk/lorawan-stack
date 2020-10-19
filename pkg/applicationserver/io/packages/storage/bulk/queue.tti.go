// Copyright © 2020 The Things Industries B.V.

package bulk

import (
	"sync"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
)

// Queue stores upstream messages.
type Queue interface {
	// Push pushes contextual upstream messages. In case of overflow, returns items not pushed in the queue. An overflow is not considered to be an error.
	Push(ups []*io.ContextualApplicationUp) ([]*io.ContextualApplicationUp, error)
	// Pop pops multiple contextual upstream messages.
	Pop(max int) ([]*io.ContextualApplicationUp, error)
}

type memoryQueue struct {
	upsMu   sync.Mutex
	ups     []*io.ContextualApplicationUp
	maxSize int
}

func (q *memoryQueue) Push(ups []*io.ContextualApplicationUp) (remaining []*io.ContextualApplicationUp, err error) {
	q.upsMu.Lock()
	defer q.upsMu.Unlock()

	count := len(ups)
	if q.maxSize > 0 && len(q.ups)+count > q.maxSize {
		count = q.maxSize - len(q.ups)
		remaining = ups[count:]
	}

	q.ups = append(q.ups, ups[:count]...)
	return remaining, nil
}

func (q *memoryQueue) Pop(count int) ([]*io.ContextualApplicationUp, error) {
	q.upsMu.Lock()
	defer q.upsMu.Unlock()
	if count < 0 || count > len(q.ups) {
		count = len(q.ups)
	}
	result := q.ups[:count]
	q.ups = q.ups[count:]
	return result, nil
}

// NewMemoryQueue instantiates a new in-memory queue.
func NewMemoryQueue(maxSize int) (Queue, error) {
	return &memoryQueue{maxSize: maxSize}, nil
}
