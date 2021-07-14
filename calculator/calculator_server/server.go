package main

import (
	"context"
	"fmt"
	pb "go-protobuf/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/status"
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
	fmt.Printf("Multiply function was invoke %v", req)
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

func (*server) Average(stream pb.CalculatorService_AverageServer) error {
	fmt.Printf("Average function was invoke ")
	sum := int32(0)
	count := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AverageRespone{
				AverageResult: float64(sum) / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("Error while reding package: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (*server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootRespone, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &pb.SquareRootRespone{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
