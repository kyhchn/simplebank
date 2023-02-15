package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kyhchn/simplebank/util"
	_ "github.com/lib/pq"
)

var testQuery *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("failed to connected to db", err)
	}
	testQuery = New(testDB)
	os.Exit(m.Run())
}
