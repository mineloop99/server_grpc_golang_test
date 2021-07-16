package main

import (
	"context"
	"fmt"
	pb "go-protobuf/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)

	//Create Blog
	fmt.Println("Creating the blog")
	blog := &pb.Blog{
		AuthorId: "WanatabeYuu",
		Title:    "2 Blog",
		Content:  "Content of the first blog",
	}
	creatBlogRes, err := c.CreateBlog(context.Background(), &pb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error when create log: %v", err)
	}

	fmt.Printf("Blog has been create: %v", creatBlogRes)
	blogId := creatBlogRes.GetBlog().GetId()

	//Read Blog

	fmt.Println("Reading the blog")
	_, err2 := c.ReadBlog(context.Background(), &pb.ReadBlogRequest{
		BlogId: "60f0fca451dcf518d2668311",
	})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}

	readBlogReq := &pb.ReadBlogRequest{BlogId: blogId}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happended while reading: %v \n", readBlogErr)
	}
	fmt.Printf("Blog has read: %v", readBlogRes)

	//update Blog

	fmt.Println("\nUpdate the blog")
	newBlog := &pb.Blog{
		Id:       blogId,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Updated Content",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &pb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was read, %v \n", updateRes)

	//delete Blog

	fmt.Println("\nDelete the blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &pb.DeleteBlogRequest{BlogId: blogId})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	//list Blog
	stream, err := c.ListBlog(context.Background(), &pb.ListBlogRequest{})

	fmt.Println("\nList Blog")
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Respone from GreetManyTime: %v", res.GetBlog())
	}
}
