package repository

import (
	"context"
	"strings"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	*mongo.Collection
}

func NewDB(db *mongo.Collection) *DB {
	return &DB{db}
}

func (db *DB) Load(key string) (model.Price, bool) {
	ticker, date := splitKey(key)
	filter := bson.D{{Key: "ticker", Value: ticker}, {Key: "price_date", Value: date}}

	var result model.Price
	err := db.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return model.Price{}, false
	}
	return model.Price{
		Open:  result.Open,
		High:  result.High,
		Low:   result.Low,
		Close: result.Close,
	}, true
}

func (db *DB) Store(key string, value model.Price) {
	ticker, date := splitKey(key)

	db.InsertOne(context.TODO(), bson.D{
		{Key: "ticker", Value: ticker},
		{Key: "price_date", Value: date},
		{Key: "open", Value: value.Open},
		{Key: "high", Value: value.High},
		{Key: "low", Value: value.Low},
		{Key: "close", Value: value.Close},
	})
}

func splitKey(key string) (string, string) {
	x := strings.Split(key, "_")
	if len(x) != 2 {
		return "", ""
	}
	return x[0], x[1]
}
