package main

import (
	"context"
	pb "gRPC_demo/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

//也就是调用的这个对象的 Search 方法
type SearchService struct {}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{
		Response: r.GetRequest() + " Server",
	}, nil
}

const PORT = "9123"

func main()  {
	server := grpc.NewServer()

	// 注册服务发现
	// 也就是把需要被调用的服务端口 注册到 gRPC Server的内部注册中心
	// 可以在接收到请求时，通过内部的服务发现
	pb.RegisterSearchServiceServer(server, &SearchService{})

	// 让服务去监听端口
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}