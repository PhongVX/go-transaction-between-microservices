package app

import (
	"database/sql"
	"fmt"
	"github.com/PhongVX/micro-protos/product"
	"google.golang.org/grpc"
	"log"
	"net"
	"orchestrator/internal/config"
	"orchestrator/internal/orderx"
	"orchestrator/internal/productx"
	"orchestrator/internal/redisx"
	"orchestrator/internal/transactioncache"
	"orchestrator/internal/transactionx"

	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
	"github.com/go-redis/redis/v8"
)

type ServerI interface {
	Start() error
	Stop()
}

func New(conf *config.Config,db *sql.DB, redisClient *redis.Client) ServerI {
	//Grpc Service
	grpcServer := grpc.NewServer()
	//TransactionCache
	txCacheSrv := transactioncache.NewTransactionCacheSrv()
	//Redis Service
	redisService := redisx.NewRedisSrv(redisClient, txCacheSrv)
	//Transaction
	txSrv := transactionx.NewTransactionSrv(db, redisService, txCacheSrv)
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
		Config: conf,
	}

}

func (a *Server) Stop() {
	a.gServer.Stop()
}

func (a *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", a.Config.Server.GPort))
	log.Printf("gRPC Server listens at port %v", a.Config.Server.GPort)
	if err != nil {
		log.Fatalf("gRPC Failed to listen at port %v %s", a.Config.Server.GPort, err)
	}
	return a.gServer.Serve(lis)
}
