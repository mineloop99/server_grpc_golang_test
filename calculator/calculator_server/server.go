package main

import (
	"context"
	"fmt"
	pb "go-protobuf/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumRespone, error) {
	fmt.Printf("Received Sum RPC: %v", in)
	firstNumber := in.FirstNumber
	secondNumber := in.SecondNumber
	sum := firstNumber + secondNumber
	res := &pb.SumRespone{
		SumResult: sum,
	}
	return res, nil
}

func (*server) Multiply(req *pb.MultiplyRequest, stream pb.CalculatorService_MultiplyServer) error {
	fmt.Printf("GreetManyTimes function was invoke %v", req)
	number := req.GetNumber()
	divisor := int32(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&pb.MultiplyRespone{
				MultiplyResult: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increase to %v", divisor)
		}
	}
	return nil
}

func main() {
	println("Calculator Server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		println("Initilize Server...")
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
