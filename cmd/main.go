package main

import (
	"log"
	"os"

	"github.com/joisandresky/go-echo-mongodb-boilerplate/configs"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/server"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/pkg/mongodb"
)

func main() {
	log.Println("STARTING SERVICE")
	configPath := configs.GetConfigPath(os.Getenv("config"))
	cfg, err := configs.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	conn := mongodb.NewConnection(cfg)
	defer conn.Close()

	server := server.NewServer(conn, cfg)

	log.Fatalln(server.Run())
}
