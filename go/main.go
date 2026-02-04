package main

import (
	"log"
	"solace/orm"
	"solace/server"
)

func main() {
	ddb := orm.New("./data/cur.db")
	srv := server.NewServer(ddb, []byte("asdf"))
	log.Fatal(srv.Start(":1323"))
}
