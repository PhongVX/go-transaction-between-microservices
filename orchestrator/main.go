package main

import (
	"context"
	"flag"
	"log"

	"orchestrator/pkg/db"
	"os"
	"os/signal"
	"time"

	"orchestrator/internal/app"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	//Connect DB
	dbConn, err := db.NewDB(&db.Config{
		Driver: "postgres",
		Host: "localhost",
		Port: 5432,
		User: "postgres",
		DBName: "postgres",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}
	server := app.New(dbConn)
	go func(){
		if err := server.Start(); err != nil {
			panic(err)
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
	server.Stop()
	log.Println("shutting down...")
	os.Exit(0)
}
