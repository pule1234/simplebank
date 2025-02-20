package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pule1234/simple_bank/util"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://sj_admin:123@localhost:5431/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file")
	}

	testDB, err = sql.Open(config.DbDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db : ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
