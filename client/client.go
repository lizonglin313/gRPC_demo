package main

import (
	"context"
	pb "gRPC_demo/proto"
	"google.golang.org/grpc"
	"log"
)

const PORT = "9123"

func main()  {
	// 连接到指定端口
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	// 创建客户端对象
	client := pb.NewSearchServiceClient(conn)
	// 发送 RCP 请求
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	// 拿到响应结果
	log.Printf("resp: %s", resp.GetResponse())
}