package repository

import (
	"context"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	collection *mongo.Collection
	logger     *zerolog.Logger
}

func NewMongo(client *mongo.Client, logger *zerolog.Logger) *Mongo {
	collection := client.Database("stock-service").Collection("prices-cache")
	return &Mongo{collection: collection, logger: logger}
}

type document struct {
	Key   string       `bson:"_id"`
	Price *model.Price `bson:"price,omitempty"`
}

func (m *Mongo) Load(key string) (model.Price, bool) {
	result := m.collection.FindOne(context.TODO(), document{Key: key})
	if result.Err() != nil {
		return model.Price{}, false
	}
	doc := document{}
	if err := result.Decode(&doc); err != nil {
		m.logger.Err(err).Send()
		return model.Price{}, false
	}
	return *doc.Price, true
}

func (m *Mongo) Store(key string, price model.Price) {
	m.collection.InsertOne(context.TODO(), document{Key: key, Price: &price})
}
