package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	mainpb "github.com/sidharth-rashwana/grpc-deepdive/proto/gen"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

type server struct {
	mainpb.UnimplementedDeepdiveServer
}

func (s *server) Add(ctx context.Context, req *mainpb.AddRequest) (*mainpb.AddResponse, error) {
	sum := req.GetA() + req.GetB()
	log.Printf("Add called with a=%d, b=%d, sum=%d", req.GetA(), req.GetB(), sum)
	return &mainpb.AddResponse{Sum: sum}, nil
}

func (s *server) GenerateFibonacci(req *mainpb.FibonacciRequest, stream mainpb.Deepdive_GenerateFibonacciServer) error {
	n := req.N
	a, b := 0, 1
	for i := 0; i < int(n); i++ {
		err := stream.Send(&mainpb.FibonacciResponse{
			Number: int32(a),
		})
		if err != nil {
			return err
		}
		log.Println("Sent number : ", int32(a))
		a, b = b, a+b
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) SendNumbers(stream mainpb.Deepdive_SendNumbersServer) error {
	var sum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mainpb.NumberResponse{
				Sum: sum,
			})
		}
		if err != nil {
			return err
		}
		log.Println(req.GetNumber())
		sum += req.GetNumber()
	}
}

func (s *server) Chat(stream mainpb.Deepdive_ChatServer) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				log.Println("Client closed send stream")
				return
			}
			if err != nil {
				log.Println("Error receiving from client:", err)
				return
			}
			log.Println("Received from client:", req.GetMessage())
		}
	}()
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter message to send (type 'quit' to end this chat): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading stdin:", err)
				return
			}
			input = strings.TrimSpace(input)
			if input == "quit" {
				log.Println("Stopping chat on server side")
				return
			}
			if err := stream.Send(&mainpb.ChatMessage{Message: input}); err != nil {
				log.Println("Error sending to client:", err)
				return
			}
		}
	}()
	<-done
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Listening on port : 50051")
	grpcServer := grpc.NewServer()
	mainpb.RegisterDeepdiveServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
