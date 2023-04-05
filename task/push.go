package task

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mygochat/model"
	"mygochat/task/rpc"
)

var msgChannel []chan model.ChatForm

func init() {
	msgChannel = make([]chan model.ChatForm, 10)
}

func (t *Task) Push(msg string) {
	var userMsg model.ChatForm
	err := json.Unmarshal([]byte(msg), &userMsg)
	if err != nil {
		log.Println(err)
		return
	}
	msgChannel[rand.Int()%10] <- userMsg
}

func (t *Task) Start() {
	for i := 0; i < len(msgChannel); i++ {
		msgChannel[i] = make(chan model.ChatForm, 100)
		go t.DealTask(msgChannel[i])
	}
}

func (t *Task) DealTask(ch chan model.ChatForm) {
	select {
	case UserMsg := <-ch:
		if UserMsg.RoomId > 0 {
			rpc.PublishRoom(UserMsg)
		}
		if UserMsg.ToUserId != "" {
			fmt.Println("deal task msg", UserMsg)
			rpc.PublishPerson(UserMsg)
		}
	}
}
