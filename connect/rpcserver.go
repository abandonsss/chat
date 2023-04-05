package connect

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mygochat/proto"
	"net"
	"strconv"
)

type Server struct {
	api_logic.UnimplementedChatServiceServer
}

func (c *Connect) StartRPC() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	api_logic.RegisterChatServiceServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) PublishPerson(_ context.Context, in *api_logic.ChatRequest) (*api_logic.ChatReply, error) {
	Userid := in.ToUserId
	uid, err := strconv.Atoi(Userid)
	fmt.Println("connect's msg is ", in)
	if err != nil {
		return &api_logic.ChatReply{}, err
	}
	bucket := Contain.Buckets[uid%Contain.BucketNumber]
	ch, err := bucket.Channel(uid)
	if err != nil {
		return &api_logic.ChatReply{}, err
	}
	ch.PushMsg(Msg{Body: in.Msg, FromUsername: in.FromUsername})
	return &api_logic.ChatReply{}, nil
}

func (s *Server) PublishRoom(_ context.Context, in *api_logic.ChatRequest) (*api_logic.ChatReply, error) {
	roomId := in.RoomId
	bucket := Contain.Buckets[int(roomId)%Contain.BucketNumber]
	room, err := bucket.Room(int(roomId))
	if err != nil {
		return &api_logic.ChatReply{}, err
	}
	room.PushMsg(Msg{Body: in.Msg, FromUsername: in.FromUsername})
	return &api_logic.ChatReply{}, nil
}
