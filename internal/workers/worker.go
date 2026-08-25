package workers

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Job func(context.Context) error
type Worker struct {
	jobs    chan Job
	running atomic.Bool
	wg      sync.WaitGroup
	errors  atomic.Uint64
}

func New(size int) *Worker {
	if size < 1 {
		size = 1
	}
	return &Worker{jobs: make(chan Job, size)}
}
func (w *Worker) Start(ctx context.Context, n int) {
	if w.running.Swap(true) {
		return
	}
	for i := 0; i < n; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j := <-w.jobs:
					if j != nil {
						if j(ctx) != nil {
							w.errors.Add(1)
						}
					}
				}
			}
		}()
	}
}
func (w *Worker) Submit(j Job) bool {
	if !w.running.Load() {
		return false
	}
	select {
	case w.jobs <- j:
		return true
	default:
		return false
	}
}
func (w *Worker) Stop() {
	if !w.running.Swap(false) {
		return
	}
	close(w.jobs)
	w.wg.Wait()
}
func (w *Worker) Errors() uint64 { return w.errors.Load() }
func Retry(ctx context.Context, attempts int, delay time.Duration, fn Job) error {
	var e error
	for i := 0; i < attempts; i++ {
		if e = fn(ctx); e == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay * time.Duration(i+1)):
		}
	}
	return e
}
func Ticker(ctx context.Context, d time.Duration, fn func()) {
	t := time.NewTicker(d)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				fn()
			}
		}
	}()
}
