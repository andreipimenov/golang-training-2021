package repository

import (
	"context"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type MongoDB struct {
	*mongo.Collection
}

func NewMongoDB(db *mongo.Collection) *MongoDB {
	return &MongoDB{db}
}

func (db *MongoDB) Load(key string) (model.Price, bool) {
	ticker, date := SplitKey(key)

	var result model.Price
	err := db.FindOne(
		context.TODO(),
		bson.D{{Key: "ticker", Value: ticker}, {Key: "date", Value: date}},
	).Decode(&result)

	if err != nil {
		return model.Price{}, false
	}
	return result, true
}

func (db *MongoDB) Store(key string, value model.Price) {
	ticker, date := SplitKey(key)
	_, err := db.InsertOne(context.TODO(),
		bson.D{
			{Key: "ticker", Value: ticker},
			{Key: "date", Value: date},
			{Key: "open", Value: value.Open},
			{Key: "high", Value: value.High},
			{Key: "low", Value: value.Low},
			{Key: "close", Value: value.Close},
		})
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Err(err).Msg("store error")
	}

}
