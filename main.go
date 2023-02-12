package main

import (
	"database/sql"
	"log"

	"github.com/kyhchn/simplebank/api"
	db "github.com/kyhchn/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable")
	if err != nil {
		log.Fatal("failed to connected to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
