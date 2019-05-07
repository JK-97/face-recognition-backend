package model

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
)

// DB is a mongodb database
var DB *mongo.Database

// InitDB initalizes mongo db
func InitDB() {
	cfg := config.Config()
	dbAddr := cfg.GetString("db-addr")

	client, err := mongo.NewClient(dbAddr)
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	DB = client.Database(cfg.GetString("app-name"))
}
