package main

import (
	"context"
	"fmt"
	pb "go-protobuf/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewCalculatorServiceClient(cc)
	//	doUnary(c)
	doServerStreaming(c)
}

func doUnary(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum Unary RPC...")
	req := &pb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 14,
	}
	in, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Respone from Sum: %v", in.SumResult)

}

func doServerStreaming(c pb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a Multiply SvStream RPC...")
	req := &pb.MultiplyRequest{
		Number: 120,
	}
	reqStream, err := c.Multiply(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while loading stream: %v", err)
	}
	for {
		msg, err := reqStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while loading stream: %v", err)
		}
		log.Printf("Respone from Multiply: %v", msg.GetMultiplyResult())
	}
}
