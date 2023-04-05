package logic

import (
	"log"
	"mygochat/logic/gomoysql"
	"mygochat/logic/goredis"
)

type Logic struct {
}

func New() *Logic {
	return new(Logic)
}

func (*Logic) Run() {
	gomoysql.InitMysql()
	err := goredis.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	StartRpc()
}
