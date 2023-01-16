package master

import (
	"context"
	"log"
	"mydoki/worker"
	"sync"
	"time"
)

type EnqueueCommand struct {
	EndNode         int
	DelayTimeInSecs int
	Message         string
}

type Usecase interface {
	Enqueue(ctx context.Context, cmd EnqueueCommand) error
}

type usecase struct {
	q            []EnqueueCommand
	firstWorker  worker.Usecase
	totalWorkers int
	mu           sync.Mutex
}

var _ Usecase = &usecase{}

func NewUsecase(totalWorkers int) *usecase {
	var lastWorker worker.Usecase

	for i := totalWorkers; i >= 1; i-- {
		lastWorker = worker.NewUsecase(i, lastWorker)
	}

	u := &usecase{
		totalWorkers: totalWorkers,
		firstWorker:  lastWorker,
		q:            make([]EnqueueCommand, 0),
	}

	go u.process(context.Background())

	return u
}

func (u *usecase) Enqueue(ctx context.Context, cmd EnqueueCommand) error {
	if cmd.EndNode <= 1 || cmd.EndNode > u.totalWorkers {
		return ErrWrongEndNode
	}

	if cmd.DelayTimeInSecs < 0 {
		return ErrWrongDelay
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Master: %s got queued.", cmd.Message)
	u.q = append(u.q, cmd)
	return nil
}

func (u *usecase) process(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			u.mu.Lock()
			if len(u.q) > 0 && u.firstWorker.IsAvailable() {
				cmd := u.q[0]
				err := u.firstWorker.Enqueue(worker.EnqueueCommand{
					EndNode:         cmd.EndNode,
					DelayTimeInSecs: cmd.DelayTimeInSecs,
					Message:         cmd.Message,
				})
				if err == nil {
					u.q = u.q[1:]
					log.Printf("Master: Message %s entered the pipeline.", cmd.Message)
				} else {
					log.Printf("Master: Message %s faced error: %v.", cmd.Message, err)
				}
			}
			u.mu.Unlock()
		}
	}
}
