package main

import (
	"database/sql"
	"github.com/MsN-12/simpleBank/api"
	db "github.com/MsN-12/simpleBank/db/sqlc"
	"github.com/MsN-12/simpleBank/util"
	_ "github.com/lib/pq"
	"log"
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
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start sever: ", err)
	}
}
