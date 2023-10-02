package main

import (
	"database/sql"
	"github.com/MsN-12/simpleBank/api"
	db "github.com/MsN-12/simpleBank/db/sqlc"
	"github.com/MsN-12/simpleBank/gapi"
	"github.com/MsN-12/simpleBank/pb"
	"github.com/MsN-12/simpleBank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("cannot load config file: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start sever: ", err)
	}
}
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalln("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}
