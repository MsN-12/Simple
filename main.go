package main

import (
	"context"
	db "github.com/MsN-12/simpleBank/db/sqlc"
	"github.com/MsN-12/simpleBank/gapi"
	"github.com/MsN-12/simpleBank/mail"
	"github.com/MsN-12/simpleBank/pb"
	"github.com/MsN-12/simpleBank/util"
	"github.com/MsN-12/simpleBank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config file: ")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to database: ")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)
	go runGateWayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

}
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance: ")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up:")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server: ")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}
	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}
}
func runGateWayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server: ")
	}
	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener: ")
	}
	log.Info().Msgf("start http gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start http gateway server")
	}
}

//	func runGinServer(config util.Config, store db.Store) {
//		server, err := api.NewServer(config, store)
//		if err != nil {
//			log.Fatal().Msg("cannot create server: ")
//		}
//		err = server.Start(config.HttpServerAddress)
//		if err != nil {
//			log.Fatal().Msg("cannot start sever: ")
//		}
//	}
