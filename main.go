package main

import (
	"mygochat/api"
	"mygochat/connect"
	"mygochat/logic"
	"mygochat/task"
)

func main() {
	go api.New().Run()
	go logic.New().Run()
	go task.New().Run()
	connect.New().Run()
}
