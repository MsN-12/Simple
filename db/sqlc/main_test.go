package db

import (
	"context"
	"github.com/MsN-12/simpleBank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("cannot load config file: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
