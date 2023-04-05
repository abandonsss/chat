package connect

import "log"

type Connect struct {
}

func New() *Connect {
	return new(Connect)
}

var Contain *Container

func (c *Connect) Run() {
	Contain = NewContainer(100)
	go c.StartRPC()
	if err := c.InitWebsocket(); err != nil {
		log.Fatal("init websocket err")
	}
}
