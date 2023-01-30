package main

import (
	"context"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang/internal/data"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	db   struct {
		uri string
	}
}

type application struct {
	models data.Models
	logger *log.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "backend server port")
	flag.StringVar(&cfg.db.uri, "db uri", "mongodb://localhost:27017", "db uri")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	client, err := openDB(cfg.db.uri)
	if err != nil {
		panic(err)
	}

	defer closeDB(*client)

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	app := &application{
		models: data.NewModels(*client),
		logger: logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func openDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	return client, err
}

func closeDB(client mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
