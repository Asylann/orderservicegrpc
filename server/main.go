package main

import (
	"context"
	pb "github.com/Asylann/OrderServiceGRPC/proto"
	"github.com/Asylann/OrderServiceGRPC/server/internal/config"
	"github.com/Asylann/OrderServiceGRPC/server/internal/repository"
	"github.com/Asylann/OrderServiceGRPC/server/internal/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env variables are loaded")
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error during load cfg: %s", err.Error())
		return
	}

	err = repository.InitDBConn(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var l net.Listener
	if l, err = net.Listen("tcp", ":"+cfg.Port); err != nil {
		log.Fatal("Can not run TCP", err.Error())
		return
	}

	orderStore, err := repository.NewOrderStore()
	if err != nil {
		log.Fatalf("Error during init orderStore:%s", err.Error())
		return
	}

	repository.InitCartServiceConn()

	srv := &service.Server{OrderStore: orderStore}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, srv)

	quit, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("Server is running on :%s \n", cfg.Port)
		if err = s.Serve(l); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-quit.Done()
	log.Println("Shut down processing...")

	done := make(chan interface{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Server is stopped running !")
	case <-time.After(10 * time.Second):
		log.Println("Server is stopped running due to timeout !")
		s.Stop()
	}
}
