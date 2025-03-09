package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/iamYole/common"
	"github.com/iamYole/common/broker"
	"github.com/iamYole/common/discovery"
	"github.com/iamYole/common/discovery/consul"
	"github.com/iamYole/oms-payments/processor/stripeprocessor"
	"github.com/stripe/stripe-go/v81"
	"google.golang.org/grpc"
)

var (
	serviceName         = "payment"
	grpcAddr            = common.GetString("GRPC_ADDR", "localhost:2002")
	consulAddr          = common.GetString("CONSUL_ADDR", "localhost:8500")
	amqpUser            = common.GetString("AMQP_USER", "amqpuser")
	amqpPassword        = common.GetString("AMQP_PASSWORD", "amqppassword")
	amqpHost            = common.GetString("AMQP_HOST", "amqphost")
	amqpPort            = common.GetString("AMQP_PORT", "amqpport")
	stripeKey           = common.GetString("STRIPE_KEY", "stripe_key")
	httpAddr            = common.GetString("HTTP_ADDR", "localhost:4242")
	endpointStripSecret = common.GetString("ENDPOINT_STRIPSECRET", "whsec_....")
)

func main() {
	//Register the Consul Server
	registry, err := consul.NewRegistery(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(context.Background(), instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	//Perform continous health check on the service
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed health check")
			}
			time.Sleep(time.Second * 2)
		}
	}()
	defer registry.Deregister(context.Background(), instanceID, serviceName)

	//setup Stripe
	stripeProcessor := stripeprocessor.NewProcessor()
	stripe.Key = stripeKey

	//Connect to the RabbitMQ Broker
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		ch.Close()
		close()
	}()

	svc := NewPaymentService(stripeProcessor)
	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	//httpServer
	mux := http.NewServeMux()
	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRoutes(mux)

	go func() {
		log.Printf("Starting HTTP Server on: %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("failed to start http server")
		}
	}()

	//gRPCServer
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to dial %s. Error: %v", grpcAddr, err.Error())
		//log.Println(err)
	}
	defer lis.Close()

	//layered architecture
	// store := NewStore()
	// svc := NewService(store)
	// NewGRPCHandler(grpcServer, svc, ch)

	log.Println("gRPC Server started at ", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
