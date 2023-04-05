package task

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Task struct {
}

func New() *Task {
	return new(Task)
}

func (t *Task) Run() {
	if err := t.InitRedis(); err != nil {
		log.Fatal(err)
	}
	t.Start()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
