package main

import (
	"database/sql"
	"example-grpc/internal/database"
	"example-grpc/internal/pb"
	"example-grpc/internal/service"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDB)

	grpcServer := grpc.NewServer()

	pb.RegisterCategoryServiceServer(grpcServer, categoryService)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		panic(err)
	}

	err = grpcServer.Serve(lis)

	if err != nil {
		panic(err)
	}
}
