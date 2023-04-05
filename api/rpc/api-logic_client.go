package rpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"mygochat/model"
	"mygochat/proto"
)

//var cfg = make(map[string]string, 5)
//
//func init() {
//	vip := viper.New()
//	vip.SetConfigName("grpc")
//	vip.SetConfigType("yaml")
//	vip.AddConfigPath("./config")
//	err := vip.ReadInConfig()
//	if err != nil {
//		log.Fatal(err)
//	}
//	cfg := vip.GetStringMapString("api-logic")
//}

func Register(request model.RegisterForm) (*model.RegisterReply, error) {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := api_logic.NewUserServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := c.Register(ctx, &api_logic.UserRegister{Username: request.UserName, Password: request.Password})
	if err != nil {
		return nil, err
	}
	rep := &model.RegisterReply{
		Msg: r.Msg,
	}
	return rep, nil
}

func Login(req model.LoginForm) (*model.LoginReply, error) {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := api_logic.NewUserServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := c.Login(ctx, &api_logic.UserLogin{Username: req.Username, Password: req.Password})
	if err != nil {
		return nil, err
	}
	reply := &model.LoginReply{Token: r.Token, Msg: r.Msg}
	return reply, nil
}

func CheckAuth(token string) (string, string, error) {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := api_logic.NewUserServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := c.CheckAuth(ctx, &api_logic.Token{Token: token})
	if err != nil {
		return "", "", err
	}
	return r.Username, r.Userid, nil
}

func Chat(chatForm model.ChatForm) error {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := api_logic.NewUserServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = c.Chat(ctx, &api_logic.ChatRequest{Msg: chatForm.Msg, FromUserId: chatForm.FromUserId, FromUsername: chatForm.FromUserName, ToUserId: chatForm.ToUserId, RoomId: chatForm.RoomId})
	return err
}
