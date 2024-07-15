// Package workerpool пул воркеров
package workerpool

import (
	"log/slog"
	"sync"
)

// WorkerPool пул воркеров, задачи передаются в канале анонимными функциями, параметры через замыкания
type WorkerPool struct {
	workerCount uint64
	wg          *sync.WaitGroup
	stopSignal  chan struct{}
	taskQueue   chan func()
}

// New создает и запускает новый воркер пул
func New(workerCount uint64) *WorkerPool {

	if workerCount < 1 {
		workerCount = 1
	}

	pool := &WorkerPool{
		workerCount: workerCount,
		wg:          &sync.WaitGroup{},
		stopSignal:  make(chan struct{}, 1),         // канал для остановки пула
		taskQueue:   make(chan func(), workerCount), // буферизованный канал для задач.
	}

	for i := uint64(0); i < workerCount; i++ {
		pool.wg.Add(1)

		go func() {
			defer pool.wg.Done()
			for {
				select {
				case task := <-pool.taskQueue:
					task()
				case <-pool.stopSignal:
					return
				}
			}
		}()
	}
	slog.Debug("WorkerPool created")
	return pool
}

// Submit добавляет задачу в пул
func (p *WorkerPool) Submit(task func()) {
	p.taskQueue <- task
}

// StopWait останавливает пул и ждет, пока все задачи будут выполнены
func (p *WorkerPool) StopWait() {
	close(p.stopSignal)
	p.wg.Wait()
	slog.Debug("WorkerPool stopped")
}
