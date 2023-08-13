package db

import (
	"database/sql"
	"github.com/MsN-12/simpleBank/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("cannot load config file: ", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
