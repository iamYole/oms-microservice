package main

import (
	"context"
	"log"
	"net"

	"github.com/iamYole/common"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.GetString("GRPC_ADDR", "localhost:2000")
)

func main() {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to dial %s. Error: %v", grpcAddr, err.Error())
		//log.Println(err)
	}
	defer lis.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc)

	svc.CreateOrder(context.Background())

	log.Println("gRPC Server started at ", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
