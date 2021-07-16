package main

import (
	"context"
	"fmt"
	pb "go-protobuf/phonebook/phonebookpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (*server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.LongGreetRespone{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while  reading client stream: %v", req)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client strean: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		sendErr := stream.Send(&pb.GreetEveryoneRespone{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}
}

func (*server) GreetWithDeadLine(ctx context.Context, in *pb.GreetWithDeadlineRequest) (*pb.GreetWithDeadlineRespone, error) {
	fmt.Println("GreetWithDeadLine func invoked %V", in)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled ")
			return nil, status.Error(codes.DeadlineExceeded, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := in.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &pb.GreetWithDeadlineRespone{
		Result: result,
	}
	return res, nil
}

func main() {
	tls := true
	var opts grpc.ServerOption
	if tls {
		certFile := "openssl/server.crt"
		keyFile := "openssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
		}
		//creds.ServerHandshake(&tls.Config{InsecureSkipVerify: true})
		opts = grpc.Creds(creds)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		println("Initilize Server...")
	}
	s := grpc.NewServer(opts)
	reflection.Register(s)
	pb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
