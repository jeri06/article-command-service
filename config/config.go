package config

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Application struct {
		Port           string
		Name           string
		AllowedOrigins []string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Mongodb struct {
		ClientOptions *options.ClientOptions
		Database      string
	}
	SaramaKafka struct {
		Addresses []string
		Config    *sarama.Config
	}
}

func Load() *Config {
	cfg := new(Config)
	cfg.logFormatter()
	cfg.app()
	cfg.sarama()
	cfg.mongodb()

	return cfg
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}
func (cfg *Config) app() {
	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	rawAllowedOrigins := strings.Trim(os.Getenv("ALLOWED_ORIGINS"), " ")

	allowedOrigins := make([]string, 0)
	if rawAllowedOrigins == "" {
		allowedOrigins = append(allowedOrigins, "*")
	} else {
		allowedOrigins = strings.Split(rawAllowedOrigins, ",")
	}

	cfg.Application.Port = port
	cfg.Application.Name = appName
	cfg.Application.AllowedOrigins = allowedOrigins
}

func (cfg *Config) sarama() {
	brokers := os.Getenv("KAFKA_BROKERS")
	sslEnable, _ := strconv.ParseBool(os.Getenv("KAFKA_SSL_ENABLE"))
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	sc := sarama.NewConfig()
	sc.Version = sarama.V2_1_0_0
	if username != "" {
		sc.Net.SASL.User = username
		sc.Net.SASL.Password = password
		sc.Net.SASL.Enable = true
	}
	sc.Net.TLS.Enable = sslEnable
	sc.Producer.Return.Successes = true
	sc.Producer.Return.Errors = true

	// consumer config
	sc.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	sc.Consumer.Offsets.Initial = sarama.OffsetOldest

	// producer config
	sc.Producer.Retry.Backoff = time.Millisecond * 500

	cfg.SaramaKafka.Addresses = strings.Split(brokers, ",")
	cfg.SaramaKafka.Config = sc
}

func (cfg *Config) mongodb() {
	appName := os.Getenv("APP_NAME")
	uri := os.Getenv("MONGODB_URL")
	db := os.Getenv("MONGODB_DATABASE")
	minPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MIN_POOL_SIZE"), 10, 64)
	maxPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MAX_POOL_SIZE"), 10, 64)
	maxConnIdleTime, _ := strconv.ParseInt(os.Getenv("MONGODB_MAX_IDLE_CONNECTION_TIME_MS"), 10, 64)
	connTimeoutMS, _ := strconv.ParseInt(os.Getenv("MONGODB_CONNECTION_TIMEOUT_MS"), 10, 64)

	opts := options.Client().
		ApplyURI(uri).
		SetAppName(appName).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(time.Millisecond * time.Duration(maxConnIdleTime)).
		SetConnectTimeout(time.Millisecond * time.Duration(connTimeoutMS))

	cfg.Mongodb.ClientOptions = opts
	cfg.Mongodb.Database = db
}
