package main

import (
	"context"
	"fmt"
	pb "go-protobuf/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doServerStreaming(c)
	//doClientStreaming(c)

	doErrorUnary(c)
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
	fmt.Printf("Starting to do a Multiply Sv Stream RPC...")
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

func doClientStreaming(c pb.CalculatorServiceClient) {
	fmt.Printf("Starting to do Average Client Stream Rpc...")
	s := [4]int32{12, 22, 23, 14}
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error while reading Context")
	}
	for _, value := range s {
		fmt.Printf("sending request: %v\n", value)
		req := &pb.AverageRequest{
			Number: value,
		}

		if err != nil {
			log.Fatalf("error while sending request: %v\n", err)
		}
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receive respone: %v", err)
	}
	fmt.Printf("respone from Server: %v\n", res)
}

func doErrorUnary(c pb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary Rpc...")

	//correct call
	doErrorCall(c, 10)

	//error call
	doErrorCall(c, -22)
}

func doErrorCall(c pb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &pb.SquareRootRequest{Number: n})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
			return
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())

}
