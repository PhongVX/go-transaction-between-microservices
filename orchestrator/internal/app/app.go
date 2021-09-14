package app

import (
	"database/sql"
	"github.com/PhongVX/micro-protos/product"
	"google.golang.org/grpc"
	"log"
	"net"
	"orchestrator/internal/orderx"
	"orchestrator/internal/productx"
	"orchestrator/internal/transactionx"

	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
)

const (
	gRPCAddress = "0.0.0.0:9999"
)

type ServerI interface {
	Start() error
	Stop()
}

func New(db *sql.DB) ServerI {
	//Grpc Service
	grpcServer := grpc.NewServer()
	//Transaction
	txSrv := transactionx.NewTransactionSrv(db)
	txGSrv := transactionx.NewGService(txSrv)
	transaction.RegisterTransactionServer(grpcServer, txGSrv)
	//Order
	orderRepos := orderx.NewRepository(txSrv)
	orderGSrv := orderx.NewGService(orderRepos)
	order.RegisterOrderServer(grpcServer, orderGSrv)
	//Product
	productRepos := productx.NewRepository(txSrv)
	productGSrv := productx.NewGService(productRepos)
	product.RegisterProductServer(grpcServer, productGSrv)

	return &Server{
		gServer: grpcServer,
	}

}

func (a *Server) Stop() {
	a.gServer.Stop()
}

func (a *Server) Start() error {
	lis, err := net.Listen("tcp", gRPCAddress)
	log.Printf("gRPC Server listens at %s", gRPCAddress)
	if err != nil {
		log.Fatalf("Failed to listen at %s %s", gRPCAddress, err)
	}
	return a.gServer.Serve(lis)
}
