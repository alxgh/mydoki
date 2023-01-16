package worker

import (
	"log"
	"time"
)

type EnqueueCommand struct {
	EndNode         int
	DelayTimeInSecs int
	Message         string
}

type Usecase interface {
	IsAvailable() bool
	Enqueue(cmd EnqueueCommand) error
}

type usecase struct {
	number      int
	nextWorker  Usecase
	isAvailable bool
	proc        chan EnqueueCommand
}

var _ Usecase = &usecase{}

func NewUsecase(number int, nextWorker Usecase) *usecase {
	u := &usecase{
		number:      number,
		nextWorker:  nextWorker,
		proc:        make(chan EnqueueCommand),
		isAvailable: true,
	}

	go u.process()

	return u
}

func (u *usecase) IsAvailable() bool {
	return u.isAvailable
}

func (u *usecase) Enqueue(cmd EnqueueCommand) error {
	if !u.IsAvailable() {
		return ErrNotAvailable
	}

	// NOTE: not concurrent safe :D
	u.isAvailable = false

	log.Printf("Worker #%d: Starting to process %s", u.number, cmd.Message)
	u.proc <- cmd

	return nil
}

func (u *usecase) process() {
	for cmd := range u.proc {
		ticker := time.NewTicker(1 * time.Second)
		iter := 0
		for {
			select {
			case <-ticker.C:
				log.Printf("Worker #%d: Processing %s, iteration: %d", u.number, cmd.Message, iter)
				iter++
			}
			if iter == cmd.DelayTimeInSecs {
				break
			}
		}
		if cmd.EndNode == u.number {
			log.Printf("Worker #%d: Finished processing %s", u.number, cmd.Message)
			u.isAvailable = true
			continue
		}
		for {
			if u.nextWorker.IsAvailable() {
				if err := u.nextWorker.Enqueue(cmd); err != nil {
					log.Printf("Worker #%d: Err enqueing message %s in next worker", u.number, cmd.Message)
				} else {
					break
				}
			}

			time.Sleep(1 * time.Second)
		}
		u.isAvailable = true
	}
}
