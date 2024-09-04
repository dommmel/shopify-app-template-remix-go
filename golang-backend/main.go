package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"

	pb "github.com/dommmel/shopify-app-template-remix-go/generated/proto"
	"github.com/dommmel/shopify-app-template-remix-go/user"
)

func main() {
	db, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, myshopifydomain TEXT, accesstoken TEXT, scopes TEXT, UNIQUE(myshopifydomain))")
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	// Setup repository and handler
	userRepo := user.NewUserRepository(db)
	userHandler := user.NewUserHandler(userRepo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userHandler)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down GRPC server...")
		s.GracefulStop()
	}()

	log.Println("Starting GRPC server on :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
