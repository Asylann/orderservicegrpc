package main

import (
	"context"
	pb "github.com/Asylann/OrderService/proto"
	"github.com/Asylann/OrderService/server/internal/config"
	"github.com/Asylann/OrderService/server/internal/models"
	"github.com/Asylann/OrderService/server/internal/repository"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"
)

func main() {
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
		log.Fatal("Can not run TCP")
		return
	}

	srv := models.Server{}

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
