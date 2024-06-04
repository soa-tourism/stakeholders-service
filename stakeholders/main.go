package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"stakeholders/handler"
	"stakeholders/model"
	"stakeholders/proto/auth"
	"stakeholders/repo"
	"stakeholders/service"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {

	dsn := "user=postgres password=super dbname=soa-stakeholders host=stakeholders_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}

	err = database.AutoMigrate(
		&model.User{},
		&model.Person{},
	)
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}

	return database
}

func startListener() net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:8089")
	if err != nil {
		log.Fatalln(err)
	}
	return listener
}

func shutdown(grpcServer *grpc.Server) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a termination signal is received
	<-stopCh

	// Shutdown gRPC server gracefully
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server gracefully stopped")
}

func main() {
	database := initDB()
	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}

	listener := startListener()
	defer listener.Close()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	userRepo := &repo.UserRepository{DB: database}
	authService := &service.AuthService{UserRepo: userRepo, Key: "i_m_done"}
	authHandler := &handler.AuthHandler{AuthService: authService}
	auth.RegisterAuthorizeServer(grpcServer, authHandler)

	reflection.Register(grpcServer)
	grpcServer.Serve(listener)

	shutdown(grpcServer)
}
