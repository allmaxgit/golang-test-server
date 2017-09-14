package main

import (
	"scoutiq_server/server"
	"scoutiq_server/db"
)

func main() {
	database := db.DatabaseConnect()
	config := &server.Config{
		SecretKey: "secret",
		ServerName: "Scoutiq REST server",
	}
	s := server.NewServer(database, config)
	s.Logger.Fatal(s.Start(":8000"))
}