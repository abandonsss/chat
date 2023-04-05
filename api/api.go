package api

import (
	"github.com/gin-gonic/gin"
	"mygochat/api/router"
)

type Chat struct {
}

func New() *Chat {
	return &Chat{}
}

func (*Chat) Run() {
	gin.SetMode(gin.ReleaseMode)
	r := router.Register()
	r.Run(":8080")
}
