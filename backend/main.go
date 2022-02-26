package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/jpittis/envoy-client-sim/backend/proto"
)

func main() {
	go func() {
		err := listen("127.0.0.1:10081", "10081")
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := listen("127.0.0.1:10082", "10082")
		if err != nil {
			log.Fatal(err)
		}
	}()

	client, conn, err := connect("127.0.0.1:10080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for {
		rep, err := client.Get(context.Background(), &pb.GetRequest{})
		if err != nil {
			log.Println("Failure!")
		} else {
			log.Printf("Success! (%s)", rep.Name)
		}
		time.Sleep(1 * time.Second)
	}
}

func listen(addr, name string) error {
	li, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBackendServer(grpcServer, &backend{name})
	return grpcServer.Serve(li)
}

type backend struct {
	name string
}

func (b *backend) Get(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {
	return &pb.GetResponse{Name: b.name}, nil
}

func connect(addr string) (pb.BackendClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return pb.NewBackendClient(conn), conn, nil
}
