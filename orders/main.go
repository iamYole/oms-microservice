package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/iamYole/common"
	"github.com/iamYole/common/broker"
	"github.com/iamYole/common/discovery"
	"github.com/iamYole/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName  = "orders"
	grpcAddr     = common.GetString("GRPC_ADDR", "localhost:2000")
	consulAddr   = common.GetString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.GetString("AMQP_USER", "amqpuser")
	amqpPassword = common.GetString("AMQP_PASSWORD", "amqppassword")
	amqpHost     = common.GetString("AMQP_HOST", "amqphost")
	amqpPort     = common.GetString("AMQP_PORT", "amqpport")
)

func main() {

	registry, err := consul.NewRegistery(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(context.Background(), instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed health check")
			}
			time.Sleep(time.Second * 2)
		}
	}()
	defer registry.Deregister(context.Background(), instanceID, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func ()  {
		ch.Close()
		close()
	}()

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to dial %s. Error: %v", grpcAddr, err.Error())
		//log.Println(err)
	}
	defer lis.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc, ch)

	//svc.CreateOrder(context.Background(),)

	log.Println("gRPC Server started at ", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
