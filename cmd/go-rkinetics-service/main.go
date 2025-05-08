package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ntquang98/go-rkinetics-service/internal/conf"
	"github.com/ntquang98/go-rkinetics-service/internal/pkg/logger"
	"github.com/spf13/viper"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs/config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func getConfig() (*conf.Bootstrap, error) {
	viper.SetConfigFile(flagconf)
	// Set default values for the configuration parameters
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.log_level", "debug")
	// set config and read
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	// merge with env config
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Explicitly bind environment variables
	viper.BindEnv("data.mongo.address", "APP_DATA_MONGO_ADDRESS")
	viper.BindEnv("data.mongo.username", "APP_DATA_MONGO_USERNAME")
	viper.BindEnv("data.mongo.password", "APP_DATA_MONGO_PASSWORD")
	viper.BindEnv("data.mongo.auth_db", "APP_DATA_MONGO_AUTH_DB")
	viper.BindEnv("data.mongo.db_name", "APP_DATA_MONGO_DB_NAME")
	viper.BindEnv("data.s3.access_key", "APP_DATA_S3_ACCESS_KEY")
	viper.BindEnv("data.s3.secret_key", "APP_DATA_S3_SECRET_KEY")
	viper.BindEnv("data.s3.region", "APP_DATA_S3_REGION")
	viper.BindEnv("data.s3.bucket", "APP_DATA_S3_BUCKET")

	// Log all Viper keys for debugging
	for _, key := range viper.AllKeys() {
		fmt.Printf("Viper key: %s, value: %v", key, viper.Get(key))
	}

	// override with ENV vars
	// for _, key := range viper.AllKeys() {
	// 	viper.Set(key, viper.Get(key))
	// }

	// unmarshal to config object
	var cfg conf.Bootstrap
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	flag.Parse()

	logger := log.With(logger.NewZapLoggerWrapper(),
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	// logger := log.With(log.NewStdLogger(os.Stdout),
	// 	"ts", log.DefaultTimestamp,
	// 	"caller", log.DefaultCaller,
	// 	"service.id", id,
	// 	"service.name", Name,
	// 	"service.version", Version,
	// 	"trace.id", tracing.TraceID(),
	// 	"span.id", tracing.SpanID(),
	// )

	// c := config.New(
	// 	config.WithSource(
	// 		file.NewSource(flagconf),
	// 	),
	// )
	// defer c.Close()

	// if err := c.Load(); err != nil {
	// 	panic(err)
	// }

	bc, err := getConfig()
	if err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
