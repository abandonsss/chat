package rpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mygochat/model"
	"mygochat/proto"
)

func CheckAuth(token string) (*model.TokenReply, error) {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &model.TokenReply{Username: "", UserId: ""}, err
	}
	c := api_logic.NewUserServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp, err := c.CheckAuth(ctx, &api_logic.Token{Token: token})
	if err != nil {
		return &model.TokenReply{Username: "", UserId: ""}, err
	}
	return &model.TokenReply{Username: resp.Username, UserId: resp.Userid}, nil
}
