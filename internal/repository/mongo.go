package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Client struct {
	*mongo.Client
	collection mongo.Collection
}

type mongoData struct {
	Ticker string
	Date   string
	Open   string
	High   string
	Low    string
	Close  string
}

type Repository interface {
	Load(string) (model.Price, bool)
	Store(string, model.Price)
}

func NewClient(mc *mongo.Client) *Client {
	tmp_collection := mc.Database("db").Collection("price")
	return &Client{mc, *tmp_collection}
}

func (client *Client) Load(key string) (model.Price, bool) {
	ticker, date := splitKey(key)
	var md mongoData
	filter := bson.D{
		{Key: "ticker", Value: ticker},
		{Key: "date", Value: date},
	}
	err := client.collection.FindOne(context.TODO(), filter).Decode(&md)
	if err != nil {
		return model.Price{}, false
	}
	return model.Price{
		Open:  md.Open,
		High:  md.High,
		Low:   md.Low,
		Close: md.Close,
	}, true
}

func (client *Client) Store(key string, value model.Price) {
	ticker, date := splitKey(key)
	md := &mongoData{ticker, date, value.Open, value.High, value.Low, value.Close}
	_, err := client.collection.InsertOne(context.TODO(), md)
	if err != nil {
		return
	}
}
