package rpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"mygochat/model"
	"mygochat/proto"
)

func PublishPerson(userMsg model.ChatForm) {
	conn, err := grpc.Dial(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	c := api_logic.NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := &api_logic.ChatRequest{
		Msg:          userMsg.Msg,
		FromUserId:   userMsg.FromUserId,
		FromUsername: userMsg.FromUserName,
		ToUserId:     userMsg.ToUserId,
	}
	_, err = c.PublishPerson(ctx, in)
	if err != nil {
		log.Println(err)
	}
}

func PublishRoom(userMsg model.ChatForm) {
	conn, err := grpc.Dial(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	c := api_logic.NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := &api_logic.ChatRequest{
		Msg:          userMsg.Msg,
		FromUserId:   userMsg.FromUserId,
		FromUsername: userMsg.FromUserName,
		RoomId:       userMsg.RoomId,
	}
	_, err = c.PublishRoom(ctx, in)
	if err != nil {
		log.Println(err)
	}
}
