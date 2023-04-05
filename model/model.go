package model

type RegisterForm struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegisterReply struct {
	Msg string
}

type LoginForm struct {
	Username string
	Password string
}

type LoginReply struct {
	Token string
	Msg   string
}

type TokenReply struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type ChatForm struct {
	ToUserId     string `json:"ToUserId"`
	Msg          string `json:"Msg"`
	FromUserName string `json:"FromUserName"`
	FromUserId   string `json:"FromUserId"`
	RoomId       int64  `json:"RoomId"`
}

type SessionForm struct {
	Token string `json:"token"`
}
