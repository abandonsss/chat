package connect

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Channel struct {
	Conn     *websocket.Conn
	UserId   int
	ChatChan chan Msg
	Lock     sync.RWMutex
	Next     *Channel
	Prev     *Channel
}

type Msg struct {
	Body         string `json:"body"`
	FromUsername string `json:"fromUsername"`
}

func NewChannel() *Channel {
	ch := new(Channel)
	ch.ChatChan = make(chan Msg, 100)
	return ch
}

func (c *Channel) PushMsg(msg Msg) {
	c.ChatChan <- msg
}

type Room struct {
	Id           int
	Lock         sync.RWMutex
	OnlineNumber int
	Next         *Channel
}

func (r *Room) PushMsg(msg Msg) {
	r.Lock.RLock()
	next := r.Next
	for next != nil {
		next.ChatChan <- msg
		next = next.Next
	}
	r.Lock.RUnlock()
}

func (r *Room) Put(ch *Channel) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	if r.Next != nil {
		r.Next.Prev = ch
	}
	ch.Next = r.Next
	r.Next = ch
	ch.Prev = nil
	r.OnlineNumber++
	return nil
}

func (r *Room) Delete() {
	r.Lock.Lock()
	next := r.Next
	r.Next = nil
	for next != nil {
		next.Prev = nil
		next2 := next.Next
		next.Next = nil
		next = next2
	}
	r.Lock.Unlock()
}
