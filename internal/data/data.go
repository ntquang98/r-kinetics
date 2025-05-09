package data

import (
	"context"

	"github.com/ntquang98/go-rkinetics-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewS3FileRepo,
	NewAnalyticsJobRepo,
)

// Data .
type Data struct {
	mongo MongoClient
	log   *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)
	mongo := NewMongoClient(c)

	if err := mongo.Connect(context.Background()); err != nil {
		log.Error("failed to connect to MongoDB: ", err)
		return nil, nil, err
	}

	cleanup := func() {
		if mongo.IsConnected() {
			if err := mongo.Database().Client().Disconnect(context.Background()); err != nil {
				log.Error("failed to disconnect MongoDB: ", err)
			}
		}
	}

	return &Data{mongo: mongo, log: log}, cleanup, nil
}

func (d *Data) Mongo() MongoClient {
	return d.mongo
}
