package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"order/internal/app"
	"order/pkg/http/request"
	"order/pkg/http/response"
	"os"
	"os/signal"
	"time"
)

type (
	Header struct {
		CorrelationID string `json:"correlationID"`
	}

	Order struct {
		ID         *string `json:"id"`
		CustomerID string  `json:"customerID"`
		TotalPrice string  `json:"totalPrice"`
	}

	OrderRequest struct {
		Header Header `json:"header"`
		Body   Order  `json:"body"`
	}
)


func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var o OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	res := &response.Base{}
	err := request.Get(fmt.Sprintf(`http://localhost:8080/api/v1/sql/transaction/begin/%s`, o.Header.CorrelationID), &res)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	id, err := InsertOrderService(o)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	//Commit
	err = request.Get(fmt.Sprintf(`http://localhost:8080/api/v1/sql/transaction/commit/%s`, o.Header.CorrelationID), &res)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, response.Base{
		ID: id,
	})
}

func InsertOrderService(o OrderRequest) (*string, error) {
	//res, err := request.Post("http://localhost:8080/api/v1/order/create", nil)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Run our server in a goroutine so that it doesn't block.
	server := app.New()
	go func() {
		if err := server.Start(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.StopGRPC()
	server.Stop(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
