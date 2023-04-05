package connect

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"mygochat/connect/rpc"
	"net/http"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Connect) InitWebsocket() error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c.websocketServer(w, r)
	})
	err := http.ListenAndServe("localhost:8080", nil)
	log.Println(err)
	return err
}

func (c *Connect) websocketServer(w http.ResponseWriter, r *http.Request) {
	upGrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade err")
		return
	}
	ch := NewChannel()
	ch.Conn = conn
	go WriteMessage(ch)
	go ReadMessage(ch)
}

func WriteMessage(ch *Channel) {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		ch.Conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			ch.Conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
			if err := ch.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		case message, ok := <-ch.ChatChan:
			ch.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				ch.Conn.WriteMessage(websocket.CloseMessage, nil)
				log.Println("chan is already close")
				return
			}
			w, err := ch.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			var b []byte
			b, err = json.Marshal(message)
			if err != nil {
				log.Println(err)
			}
			w.Write(b)
			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func ReadMessage(ch *Channel) {
	defer func() {
		ch.Conn.Close()
	}()

	ch.Conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	ch.Conn.SetPongHandler(func(string) error {
		ch.Conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		return nil
	})
	for {
		_, message, err := ch.Conn.ReadMessage()
		if err != nil {
			log.Println("read message err")
		}
		if message == nil {
			return
		}
		var t User
		if err := json.Unmarshal(message, &t); err != nil {
			log.Println(err)
			return
		}
		reply, err := rpc.CheckAuth(t.Token)
		if err != nil {
			log.Println(err)
			log.Println("the token is not correct")
			return
		}
		userId, err := strconv.Atoi(reply.UserId)
		if err != nil {
			log.Println(err)
			return
		}
		bucket := Contain.Buckets[userId%Contain.BucketNumber]
		err = bucket.Push(userId, t.RoomId, ch)
		if err != nil {
			log.Println(err)
			ch.Conn.Close()
			return
		}
		log.Println("connect success")
	}
}

type User struct {
	RoomId int    `json:"roomId"`
	Token  string `json:"token"`
}
