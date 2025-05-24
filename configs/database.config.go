package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (c *Config) ConnectDatabase() *Config {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(c.MongoDB_URI),
		options.Client().SetConnectTimeout(time.Second*10),
		options.Client().SetTimeout(time.Second*60))

	if err != nil {
		log.Fatalf("connect to db -> %s failed: %v", c.MongoDB_URI, err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("pinging to db -> %s failed: %v", c.MongoDB_URI, err)
	}

	log.Printf("connected to db -> %s", c.MongoDB_URI)
	c.MongoDB = client
	return c
}

func (c *Config) GetDatabase() *mongo.Database {
	return c.MongoDB.Database(c.MongoDB_Database)
}
