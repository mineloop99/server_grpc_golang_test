package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "go-protobuf/phonebook/phonebookpb"
	"io"

	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func createClientFromFile(certFile, serverNameOverride string) (credentials.TransportCredentials, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp, InsecureSkipVerify: true}), nil
}

func main() {

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "openssl/ca.crt"
		creds, sslErr := createClientFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}

		opts = grpc.WithTransportCredentials(creds)
	}
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewGreetServiceClient(cc)
	doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	//doBiDiStreaming(c)

	//doUnaryWithDeadLine(c, 5*time.Second)

	//doUnaryWithDeadLine(c, 1*time.Second)
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

func doClientStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	request := []*pb.LongGreetRequest{
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wanatabe",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wa22",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "WanatabeCC",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Hello",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range request {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving respone from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Respone: %v\n", res)
}

func doBiDiStreaming(c pb.GreetServiceClient) {
	fmt.Println("Staring to do a BiDi Streaming Rpc...")
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}
	requests := []*pb.GreetEveryoneRequest{
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wanatabe",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wa22",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "WanatabeCC",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Hello",
			},
		},
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)

			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func doUnaryWithDeadLine(c pb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline Rpc...")
	req := &pb.GreetWithDeadlineRequest{
		Greeting: &pb.Greeting{
			FirstName: "Wanatabe",
			LastName:  "Maarek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline Rpc: %v", err)
		}
		log.Fatalf("error while calling Greet RPC: %v", err)
		return
	}
	log.Printf("Respone from Greet: %v", res.Result)
}
