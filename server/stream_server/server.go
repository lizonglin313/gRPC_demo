package main

import (
	pb "gRPC_demo/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

// 下面有相应的三个方法去实现了这个接口
type StreamService struct{}

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

// 服务端流式 RPC 也就是 客户发起普通的RPC请求 服务器流式响应
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// 客户端流式
// 客户端流式请求，服务端以正常的方式响应
// 服务端持续接收请求流
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			// SendAndClose 是发现请求流关闭后
			// 将响应结果发送给客户端
			// 同时把等待接收的 Recv 给关上
			return stream.SendAndClose(&pb.StreamResponse{
				Pt: &pb.StreamPoint{
					Name:  "gRPC Stream Server: Record",
					Value: 1,
				},
			})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %s, pt.value %d", r.Pt.Name, r.Pt.Value)
	}
	return nil
}

// 双向流式
// 客户端 流式请求
// 服务端流式响应
// 首个请求一定是客户端先发
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gRPC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("stream.Recv-%d pt.name: %s, pt.value: %d", n-1, r.Pt.Name, r.Pt.Value)
	}
	return nil
}
