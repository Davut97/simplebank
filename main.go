package main

import (
	"database/sql"
	"log"

	"github.com/Davut97/simplebank/api"
	db "github.com/Davut97/simplebank/db/sqlc"
	"github.com/Davut97/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db")
	}

	store := db.NewStore(conn)
	server,err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connect start server: ", err)
	}
}
