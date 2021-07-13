package main

import (
	"context"
	"fmt"
	pb "go-protobuf/phonebook/phonebookpb"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetRespone, error) {
	fmt.Println("greet func invoked %V", in)
	firstName := in.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &pb.GreetRespone{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *pb.GreetManyTimesRequest, stream pb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoke %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number" + strconv.Itoa(i)
		res := &pb.GreetManyTimesRespone{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		println("Initilize Server...")
	}
	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
