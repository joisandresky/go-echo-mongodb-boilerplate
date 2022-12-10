package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joisandresky/go-echo-mongodb-boilerplate/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection interface {
	Close()
	DB(cfg *configs.Config) *mongo.Database
}

type mongoConnection struct {
	client *mongo.Client
}

func NewConnection(cfg *configs.Config) MongoConnection {
	var c mongoConnection
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var clientOptions *options.ClientOptions
	url := getURL(cfg.MongoDB.Host, cfg.MongoDB.DbName)

	creds := options.Credential{
		Username: cfg.MongoDB.User,
		Password: cfg.MongoDB.Password,
	}

	clientOptions = options.Client().ApplyURI(url).SetAuth(creds)
	c.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panicln("[MONGODB] Error to Connect Database", err)
	}

	err = c.client.Ping(ctx, nil)
	if err != nil {
		log.Panicln("[MONGODB] Error to Ping Database!", err)
	}

	log.Println("[MONGODB] Database Connected")

	return &c
}

func (c *mongoConnection) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.client.Disconnect(ctx)
	log.Println("[MONGODB] Disconnecting Database ...")
}

func (c *mongoConnection) DB(cfg *configs.Config) *mongo.Database {
	return c.client.Database(cfg.MongoDB.DbName)
}

func getURL(host string, dbName string) string {
	return fmt.Sprintf("mongodb://%s/%s?retryWrites=true&w=majority",
		host,
		dbName,
	)
}
