package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/product"
	"github.com/PhongVX/micro-protos/transaction"
	"log"
	"net"
	"orchestrator/internal/productx"
	"orchestrator/internal/transactionx"
	"orchestrator/pkg/db"
	"os"
	"os/signal"
	"time"

	"orchestrator/internal/orderx"

	"google.golang.org/grpc"
)

type (
	Header struct {
		CorrelationID string `json:"correlationID"`
	}

	CustomerRequest struct {
		Header Header              `json:"header"`
		Body   CustomerRequestBody `json:"body"`
	}

	CustomerRequestBody struct {
		ID             *string `json:"id"`
		LastTotalPayed string  `json:"totalPayed"`
	}
)

var tranSrv = transactionx.NewTransactionSrv(db.DBCon)


func UpdateCustomer(c CustomerRequest) (int64, error) {
	tx, err := tranSrv.GetTxByCorrelationID(c.Header.CorrelationID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	queryString := `UPDATE customer SET last_total_payed=$1 WHERE id=$2`
	rs, err := tx.Exec(queryString, c.Body.LastTotalPayed, c.Body.ID)
	if val, err := rs.RowsAffected(); err == nil && val > 0 {
		return val, nil
	}
	return 0, err
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	//Grpc Service
	lis, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatal("Failed to listen on port 9999 %s", err)
	}
	grpcServer := grpc.NewServer()
	//Transaction
	txSrv := transactionx.NewTransactionSrv(db.DBCon)
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
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	_, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	//srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
