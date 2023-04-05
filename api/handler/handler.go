package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"mygochat/api/rpc"
	"mygochat/model"
	"net/http"
)

func Register(c *gin.Context) {
	var register model.RegisterForm
	if err := c.ShouldBindBodyWith(&register, binding.JSON); err != nil {
		log.Fatal(err)
		return
	}
	reply, err := rpc.Register(register)
	if err != nil {
		log.Println(reply.Msg)
		return
	}
	c.JSON(http.StatusOK, "register success")
}

func Login(c *gin.Context) {
	var login model.LoginForm
	if err := c.ShouldBindBodyWith(&login, binding.JSON); err != nil {
		c.JSON(http.StatusOK, err)
		log.Println(err)
	}
	reply, err := rpc.Login(login)
	fmt.Println(reply.Token)
	if err != nil {
		c.JSON(http.StatusOK, err)
		log.Println(err)
	}
	c.JSON(http.StatusOK, "Login success and the token is : "+reply.Token)
}

func CheckAuth(c *gin.Context) {
	var auth model.SessionForm
	if err := c.ShouldBindBodyWith(&auth, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, "please add token")
		log.Println(err)
	}
	fmt.Println(auth.Token)
	username, userid, err := rpc.CheckAuth(auth.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "token is not correct")
	}
	c.Request.Header.Set("username", username)
	c.Request.Header.Set("userid", userid)
}

func Chat(c *gin.Context) {
	var chatForm model.ChatForm
	if err := c.ShouldBindBodyWith(&chatForm, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, "no body")
	}
	chatForm.RoomId = 0
	username := c.Request.Header.Get("username")
	if username == "" {
		c.JSON(http.StatusOK, "no username")
	}
	userid := c.Request.Header.Get("userid")
	if userid == "" {
		c.JSON(http.StatusOK, "no userid")
	}
	chatForm.FromUserId = userid
	chatForm.FromUserName = username
	fmt.Println(chatForm)
	err := rpc.Chat(chatForm)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}
	c.JSON(http.StatusOK, "send success")
}

func ChatRoom(c *gin.Context) {
	var chatForm model.ChatForm
	//if err := c.ShouldBindWith(&chatForm, binding.JSON); err != nil {
	//	c.JSON(http.StatusBadRequest, "no body")
	//}
	chatForm.ToUserId = ""
	username := c.Request.Header.Get("username")
	if username == "" {
		c.JSON(http.StatusOK, "no username")
	}
	userid := c.Request.Header.Get("userid")
	if userid == "" {
		c.JSON(http.StatusOK, "no userid")
	}
	chatForm.FromUserId = userid
	chatForm.FromUserName = username
}
