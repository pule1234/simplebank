package main

import (
	"database/sql"
	"github.com/pule1234/simple_bank/api"
	db "github.com/pule1234/simple_bank/db/sqlc"
	"github.com/pule1234/simple_bank/gapi"
	"github.com/pule1234/simple_bank/pb"
	"github.com/pule1234/simple_bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DbDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	runGrpcServer(*config, store)
}

// 不需要定义api 通过rpc方式直接调用
func runGrpcServer(
	config util.Config,
	store db.Store,
) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listen, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot listen:", err)
	}
	log.Println("gRPC server listening on", listen.Addr().String())
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
