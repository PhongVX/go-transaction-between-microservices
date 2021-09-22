package app

import (
	"context"
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"order/internal/route"
)

var (
	serviceAddress = ":8081"
	gRPCAddres = "localhost:9999"
)

func New() AppI {
	//=====================gRPC==========================
	gRPCClient, err := grpc.Dial(gRPCAddres, grpc.WithInsecure())
	if err != nil {
		log.Panic("Cound't connect to gRPC Server")
	}
	orderC := order.NewOrderClient(gRPCClient)
	txC := transaction.NewTransactionClient(gRPCClient)

	//==================Api Server====================
	r, err := route.NewRouter(orderC, txC)
	if err != nil {
		log.Panicf("Failed Init Router")
	}
	httpServer := &http.Server{
		Addr:    serviceAddress,
		Handler: r,
	}
	return &App{
		Server: httpServer,
		GRPCClient: gRPCClient,
	}
}

func (a *App) StopGRPC() error {
	return a.GRPCClient.Close()
}

func (a *App) Start() error {
	log.Printf("Server is listening at %s", serviceAddress)
	return a.Server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}