package schedular

import (
	"fmt"
	"time"
)

type Scheduler struct{}

func New() Scheduler {
	return Scheduler{}
}

// * long-ruuning process
func (s Scheduler) Start(done <-chan bool) {
	fmt.Println("scheduler started")
	for {
		select {
		case <-done:
			fmt.Println("Exiting...")
			return
		default:
			now := time.Now()
			fmt.Println("scheduler now", now)

			time.Sleep(3 * time.Second)
		}
	}
}
