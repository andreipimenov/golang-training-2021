package repository

import (
	"context"
	"fmt"
	cfg "github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const Format = "mongodb://%s:%s@%s:%s/?authSource=admin"

type MongoPrice struct {
	Id      string      `bson:"_id"`
	Price   model.Price `bson:"price"`
	Changed time.Time   `bson:"changed"`
}

type Mongo struct {
	client *mongo.Client
}

func NewMongo() (*Mongo, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf(
		Format,
		cfg.Get().MongoSettings.User,
		cfg.Get().MongoSettings.Password,
		cfg.Get().MongoSettings.Host,
		cfg.Get().MongoSettings.Port)))
	if err != nil {
		log.Fatal().Err(err)
	}
	err = client.Ping(ctx, &readpref.ReadPref{})
	if err != nil {
		log.Fatal().Err(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &Mongo{client: client}, mongoDisconnect(client)
}

func (m *Mongo) Load(key string) (model.Price, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	curs, err := m.collection().Find(ctx, bson.M{"_id": key})
	if err != nil {
		log.Debug().Err(err)
		return model.Price{}, false
	}
	defer closeCursor(curs, ctx)()
	var r []MongoPrice
	if err = curs.All(ctx, &r); err != nil {
		log.Debug().Err(err)
		return model.Price{}, false
	}
	if len(r) > 0 {
		return r[0].Price, true
	}
	return model.Price{}, false
}

func closeCursor(cur *mongo.Cursor, ctx context.Context) func() {
	return func() {
		err := cur.Close(ctx)
		if err != nil {
			log.Debug().Err(err)
		}
	}
}

func (m *Mongo) Store(key string, value model.Price) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d := []interface{}{
		MongoPrice{
			Id:      key,
			Price:   value,
			Changed: time.Now(),
		},
	}
	_, err := m.collection().InsertMany(ctx, d)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			for _, doc := range d {
				_, err := m.collection().ReplaceOne(ctx, bson.M{"_id": doc.(MongoPrice).Id}, doc)
				if err != nil {
					log.Debug().Err(err)
				}
			}
		}
		log.Debug().Err(err)
	}
}

func (m *Mongo) collection() *mongo.Collection {
	return m.client.Database("backend").Collection("prices")
}

func mongoDisconnect(client *mongo.Client) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		log.Debug().Msg("Mongo disconnect")
		err := client.Disconnect(ctx)
		if err != nil {
			log.Debug().Err(err)
		}
	}
}
