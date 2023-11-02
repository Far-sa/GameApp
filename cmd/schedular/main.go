package main

import (
	"fmt"
	"game-app/config"
	"game-app/schedular"
	"os"
	"os/signal"
	"time"
)

func main() {
	// TODO read config path from cmd
	cfg := config.Load("config.yml")
	fmt.Printf("cfg : %+v\n", cfg)

	done := make(chan bool)
	go func() {
		sch := schedular.New()
		sch.Start(done)
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	fmt.Println("received interrupt signal,shutting down gracefully...")

	done <- true
	time.Sleep(cfg.Application.GracefullShutdownTimeout)

}
