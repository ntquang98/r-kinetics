package data

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ntquang98/go-rkinetics-service/internal/conf"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoClient interface {
	Connect(context.Context) error
	Database() *mongo.Database
	IsConnected() bool
}

type Isolation struct {
	Read  *readconcern.ReadConcern
	Write *writeconcern.WriteConcern
}

type MongoConfig struct {
	Address            string
	Username           string
	Password           string
	AuthDB             string
	DBName             string
	SSL                bool
	SecondaryPreferred bool
	DoWriteTest        bool
}

type Database = mongo.Database
type SessionContext = mongo.SessionContext

type Client struct {
	config      *MongoConfig
	client      *mongo.Client
	database    *mongo.Database
	isConnected bool
}

var (
	heartbeatInterval = 60 * time.Second
	maxIdleTime       = 180 * time.Second
	minPoolSize       = uint64(2)
	defaultIsolation  = Isolation{
		Read:  readconcern.Snapshot(),
		Write: writeconcern.Majority(),
	}
)

func NewMongoClient(c *conf.Data) MongoClient {
	return &Client{
		config: &MongoConfig{
			Address:            c.Mongo.Address,
			Username:           c.Mongo.Username,
			Password:           c.Mongo.Password,
			AuthDB:             c.Mongo.Authdb,
			DBName:             c.Mongo.Dbname,
			SSL:                c.Mongo.Ssl,
			SecondaryPreferred: c.Mongo.SecondaryPreferred,
			DoWriteTest:        c.Mongo.DoWriteTest,
		},
	}
}

func (c *Client) Connect(ctx context.Context) error {
	env := os.Getenv("env")
	appName, err := os.Hostname()
	if err != nil {
		appName = "Unknown"
	}
	appName += "_" + env

	opt := &options.ClientOptions{
		AppName:           &appName,
		HeartbeatInterval: &heartbeatInterval,
		MaxConnIdleTime:   &maxIdleTime,
		MinPoolSize:       &minPoolSize,
	}
	if c.config.AuthDB != "" {
		opt.Auth = &options.Credential{
			AuthSource: c.config.AuthDB,
			Username:   c.config.Username,
			Password:   c.config.Password,
		}
	}
	if c.config.SecondaryPreferred {
		opt.ReadPreference = readpref.SecondaryPreferred()
	}
	opt.ApplyURI(c.config.Address)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	c.client = client
	c.database = client.Database(c.config.DBName)
	c.isConnected = true

	if c.config.DoWriteTest {
		inst := NewDBInstance[*ConnectionLog]("_db_connection")
		inst.ApplyDatabase(c.database)
		if _, err := inst.Create(context.TODO(), ConnectionLog{
			Host:          appName,
			ConnectedTime: time.Now(),
		}); err != nil {
			return fmt.Errorf("write test failed: %v", err)
		}
	}

	return nil
}

func (c *Client) Database() *mongo.Database {
	return c.database
}

func (c *Client) IsConnected() bool {
	return c.isConnected
}

type ConnectionLog struct {
	ID              *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedTime     *time.Time          `bson:"created_time,omitempty"`
	LastUpdatedTime *time.Time          `bson:"last_updated_time,omitempty"`
	Host            string              `bson:"host,omitempty"`
	ConnectedTime   time.Time           `bson:"connected_time,omitempty"`
}
