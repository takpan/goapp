package watcher

import (
	"sync"

	"github.com/google/uuid"
)

type Watcher struct {
	id          string        // Watcher ID.
	outCh       chan *Counter // Updates to counter will notify this channel.
	counter     *Counter      // The counter.
	counterLock *sync.RWMutex // Lock for counter.
}

func New() *Watcher {
	w := Watcher{}
	w.id = uuid.NewString()
	w.outCh = make(chan *Counter, 1)
	w.counter = &Counter{Iteration: 0}
	w.counterLock = &sync.RWMutex{}
	return &w
}

func (w *Watcher) GetWatcherId() string { return w.id }

func (w *Watcher) Send(str string) {
	w.counter.Iteration++
	w.counter.Value = str
	w.outCh <- w.counter
}

func (w *Watcher) Recv() <-chan *Counter { return w.outCh }

func (w *Watcher) ResetCounter() {
	w.counterLock.Lock()
	defer w.counterLock.Unlock()

	w.counter.Iteration = 0
	w.outCh <- w.counter
}

func (w *Watcher) Reset() {
	w.id = uuid.NewString()
	w.counter.Iteration = 0
	w.outCh = make(chan *Counter, 1)
}
