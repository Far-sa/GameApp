package schedular

import (
	"fmt"
	"game-app/param"
	"game-app/service/matchingservice"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Config struct {
	MatchWaitedUsersIntervalInSeconds int `koanf:"match_waited_users_inerval_in_seconds" `
}

type Scheduler struct {
	config   Config
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(config Config, matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		config:   config,
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

// * long-ruuning process
func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Every(s.config.MatchWaitedUsersIntervalInSeconds).Second().Do(s.MatchWaitedUsers)

	s.sch.StartAsync()

	<-done

	fmt.Println("stop schedular...")
	s.sch.Stop()

}

func (s Scheduler) MatchWaitedUsers() {
	fmt.Println("matching users...", time.Now())
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
