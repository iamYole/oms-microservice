package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/iamYole/common"
	"github.com/iamYole/common/discovery"
	"github.com/iamYole/common/discovery/consul"
	"github.com/iamYole/oms-gateway/gateway"
)

var (
	httpAddr    = common.GetString("HTTP_ADDR", ":3000")
	consulAddr  = common.GetString("CONSUL_ADDR", "localhost:8500")
	serviceName = "gateway"
)

func main() {
	registry, err := consul.NewRegistery(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(context.Background(), instanceID, serviceName, httpAddr); err != nil {
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

	mux := http.NewServeMux()
	ordersGatway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(ordersGatway)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal(("Failed to start http server"))
	}

}
