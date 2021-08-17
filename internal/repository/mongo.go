package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Data struct {
	ticker     string
	price_date string
	open       string
	high       string
	low        string
	close      string
}

type Client struct {
	collection mongo.Collection
}

func NewClient(c *mongo.Client) *Client {
	collection := c.Database("db").Collection("prices")
	return &Client{*collection}
}

func (c *Client) Load(key string) (model.Price, bool) {
	ticker, date := splitKey(key)
	var data Data
	err := c.collection.FindOne(context.TODO(), bson.D{
		{Key: ticker},
		{Value: date},
	})
	if err != nil {
		return model.Price{}, false
	}

	return model.Price{
		Open:  data.open,
		High:  data.close,
		Low:   data.low,
		Close: data.close,
	}, true
}

func (c *Client) Store(key string, value model.Price) {
	ticker, date := splitKey(key)
	data := Data{ticker: ticker, price_date: date, open: value.Open, high: value.High, low: value.Low, close: value.Close}
	c.collection.InsertOne(context.TODO(), &data)
}
