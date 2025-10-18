package main

import (
	proto "content/api/gen/content/v1"
	"content/api/internal/service"
	"content/api/internal/storage"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main(){
	path := "./data/files.db"
	stg, err := storage.NewFileStorageSQlite(path, "./data")
	if err != nil{
		log.Fatalf("Error Storage NewFileStorageSQlite: %v", err)
	}
	serv := service.NewFileService(stg)
	limiter := service.NewLimiters()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(limiter.UnaryInterceptor),
		grpc.StreamInterceptor(limiter.StreamInterceptor),
	)
	proto.RegisterContentServiceServer(server, serv)

	live, err := net.Listen("tcp", ":8080")
	if err != nil{
		log.Fatal("Error Listen: %v", err)
	}

	log.Println("Start Server 8080")
	err = server.Serve(live)
	if err != nil{
		log.Fatal("Error Server Live: %v")
	}
}