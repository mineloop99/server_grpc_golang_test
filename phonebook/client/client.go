package main

import (
	"context"
	"fmt"
	pb "go-protobuf/phonebook/phonebookpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewGreetServiceClient(cc)
	//doUnary(c)

	doServerStreaming(c)
}

func doUnary(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "Wanatabe",
			LastName:  "Yuu",
		},
	}
	in, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Respone from Greet: %v", in.Result)

}

func doServerStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &pb.GreetManyTimesRequest{
		Greeting: &pb.Greeting{
			FirstName: "Wanatabe",
			LastName:  "Yuu",
		},
	}
	reqStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling ManyTimes RPC: %v", err)
	}
	for {
		msg, err := reqStream.Recv()
		if err == io.EOF {
			// reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Respone from GreetManyTime: %v", msg.GetResult())
	}
}
