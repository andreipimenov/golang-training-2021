package repository

import (
	"context"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog"
)

type DB struct {
	*mongo.Collection
}

func NewDB(db *mongo.Collection) *DB {
	return &DB{db}
}

func (db *DB) Load(key string) (model.Price, bool) {
	ticker, date := splitKey(key)
	var result model.Price

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur := db.FindOne(ctx, bson.M{"ticker": ticker, "date": date})
	err := cur.Decode(&result)
	if err != nil {
		return model.Price{}, false
	}
	return result, true
}

func (db *DB) Store(key string, value model.Price) {
	ticker, date := splitKey(key)
	var open, high, low, close string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.InsertOne(ctx, bson.M{"ticker": ticker, "open": open,
		"high": high, "low": low,
		"close": close, "date": date})

	if err != nil {
		//creating a new logger to not edit the interface
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Err(err).Msg("function store occused an error")
	}
}

func splitKey(key string) (string, string) {
	x := strings.Split(key, "_")
	if len(x) != 2 {
		return "", ""
	}
	return x[0], x[1]
}
