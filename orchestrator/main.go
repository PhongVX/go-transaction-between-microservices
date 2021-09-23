package main

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
	"log"
	"orchestrator/pkg/config"
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
	//TODO Read from config file instead of hard code
	conf, err := config.ReadConfig("./internal/config/config.yaml")
	if err != nil {
		panic(err)
	}

	//Connect Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:       conf.Redis.Addr,
		Password:   conf.Redis.Password,
		DB:         conf.Redis.DB,
		MaxRetries: conf.Redis.MaxRetries,
	})

	//Connect DB
	dbConn, err := db.NewDB(&db.Config{
		Driver:   conf.Database.Driver,
		Host:     conf.Database.Host,
		Port:     conf.Database.Port,
		User:     conf.Database.User,
		DBName:   conf.Database.DBName,
		Password: conf.Database.Password,
	})
	if err != nil {
		panic(err)
	}

	server := app.New(conf, dbConn, redisClient)

	go func() {
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
	dbConn.Close()
	redisClient.Close()
	log.Println("shutting down...")
	os.Exit(0)
}
