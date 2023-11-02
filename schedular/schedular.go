package schedular

import (
	"fmt"
	"game-app/param"
	"game-app/service/matchingservice"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

// * long-ruuning process
func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Every(3).Second().Do(s.MatchWaitedUsers)

	s.sch.StartAsync()

	<-done

	fmt.Println("stop schedular...")
	s.sch.Stop()

}

func (s Scheduler) MatchWaitedUsers() {
	fmt.Println("matching users...", time.Now())
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
