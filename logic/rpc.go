package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mygochat/logic/gomoysql"
	"mygochat/logic/goredis"
	"mygochat/logic/tool"
	"mygochat/model"
	"mygochat/proto"
	"net"
)

type Server struct {
	api_logic.UnimplementedUserServiceServer
}

func (s *Server) Register(_ context.Context, in *api_logic.UserRegister) (*api_logic.RespRegister, error) {
	username := in.GetUsername()
	password := in.GetPassword()
	if username == "" || password == "" {
		return &api_logic.RespRegister{Msg: ""}, errors.New("username or password is empty")
	}
	user := gomoysql.QueryUsername(username)
	if user.ID > 0 {
		return &api_logic.RespRegister{Msg: ""}, errors.New("username already exist")
	}
	user.Username = username
	user.Password = password
	gomoysql.AddUser(&user)
	return &api_logic.RespRegister{Msg: "register success"}, nil
}

func (s *Server) Login(_ context.Context, in *api_logic.UserLogin) (*api_logic.RespLogin, error) {
	username := in.GetUsername()
	password := in.GetPassword()
	if username == "" || password == "" {
		return &api_logic.RespLogin{Token: "", Msg: ""}, errors.New("username or password is empty")
	}
	user := gomoysql.QueryUsername(username)
	if user.ID <= 0 {
		return &api_logic.RespLogin{Token: "", Msg: ""}, errors.New("username is not exist")
	}
	if user.Password != password {
		return &api_logic.RespLogin{Token: "", Msg: ""}, errors.New("password not correct")
	}

	sessionId := tool.GenerateToken()
	userdata := make(map[string]interface{}, 5)
	userdata["username"] = user.Username
	userdata["userid"] = user.ID
	err := goredis.HMSet(sessionId, userdata)
	if err != nil {
		return &api_logic.RespLogin{Token: "", Msg: "redis set err"}, err
	}
	return &api_logic.RespLogin{Token: sessionId, Msg: "login success"}, nil
}

func (s *Server) CheckAuth(_ context.Context, in *api_logic.Token) (*api_logic.RespToken, error) {
	sessionID := in.Token
	fmt.Println("输入token", sessionID)
	val, err := goredis.HGetALL(sessionID)
	fmt.Println(val)
	if err != nil {
		log.Println(err)
	}
	username, ok := val["username"]
	if !ok {
		return &api_logic.RespToken{Username: ""}, errors.New("not find username")
	}
	userid, ok := val["userid"]
	if !ok {
		return &api_logic.RespToken{Username: ""}, errors.New("mot find userid")
	}
	return &api_logic.RespToken{Username: username, Userid: userid}, nil
}

func (s *Server) Chat(_ context.Context, in *api_logic.ChatRequest) (*api_logic.ChatReply, error) {
	var msg model.ChatForm
	msg.FromUserId = in.FromUserId
	msg.FromUserName = in.FromUsername
	msg.RoomId = in.RoomId
	msg.ToUserId = in.ToUserId
	msg.Msg = in.Msg
	b, err := json.Marshal(msg)
	if err != nil {
		return &api_logic.ChatReply{}, err
	}
	err = goredis.LPush("chat", b)
	return &api_logic.ChatReply{}, err
}

func StartRpc() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	api_logic.RegisterUserServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
