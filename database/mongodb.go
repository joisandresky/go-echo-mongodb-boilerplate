package database

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

type Connection interface {
	Close()
	DB() *mongo.Database
}

type conn struct {
	client *mongo.Client
}

func NewConnection() Connection {
	var c conn
	var err error
	var clientOptions *options.ClientOptions
	url := getURL()

	creds := options.Credential{
		Username: viper.GetString("database.user"),
		Password: viper.GetString("database.pwd"),
	}

	clientOptions = options.Client().ApplyURI(url).SetAuth(creds)
	c.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Panicln("Error to Connect Database", err)
	}

	err = c.client.Ping(ctx, nil)
	if err != nil {
		log.Panicln("Error to Ping Database", err)
	}

	log.Println("Connected to Database!")

	return &c
}

func (c *conn) Close() {
	c.client.Disconnect(ctx)
	log.Println("Disconnecting Database!")
}

func (c *conn) DB() *mongo.Database {
	return c.client.Database(viper.GetString("database.db_name"))
}

func getURL() string {
	return fmt.Sprintf("mongodb+srv://%s/%s?retryWrites=true&w=majority",
		viper.GetString("database.host"),
		viper.GetString("database.db_name"),
	)
}
