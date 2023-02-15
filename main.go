package main

import (
	"database/sql"
	"log"

	"github.com/kyhchn/simplebank/api"
	db "github.com/kyhchn/simplebank/db/sqlc"
	"github.com/kyhchn/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("failed to connected to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
