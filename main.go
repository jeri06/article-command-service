package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jeri06/article-command-service/article"
	"github.com/jeri06/article-command-service/config"
	"github.com/jeri06/article-command-service/mongodb"
	"github.com/jeri06/article-command-service/response"
	"github.com/jeri06/article-command-service/server"
	_ "github.com/joho/godotenv/autoload" //
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cfg          *config.Config
	utcLocation  *time.Location
	indexMessage string = "Application is running properly"
)

func init() {
	utcLocation, _ = time.LoadLocation("UTC")
	cfg = config.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	vld := validator.New()

	mc, err := mongo.NewClient(cfg.Mongodb.ClientOptions)
	if err != nil {
		logger.Fatal(err)
	}
	mca := mongodb.NewClientAdapter(mc)
	if err := mca.Connect(context.Background()); err != nil {
		logger.Fatal(err)
	}
	mdb := mca.Database(cfg.Mongodb.Database)
	//kafka

	producers, err := sarama.NewSyncProducer(cfg.SaramaKafka.Addresses, cfg.SaramaKafka.Config)
	if err != nil {
		logrus.Errorf("Unable to create kafka producer got error %v", err)
		return
	}
	defer func() {
		if err := producers.Close(); err != nil {
			logrus.Errorf("Unable to stop kafka producer: %v", err)
			return
		}
	}()

	logrus.Infof("Success create kafka sync-producer")

	router := mux.NewRouter()
	router.HandleFunc("/command-service", index)

	articleRepository := article.NewArticleRepository(logger, mdb)
	articleUsecase := article.NewArticleUsecase(article.UsecaseProperty{
		ServiceName: cfg.Application.Name,
		UTCLoc:      utcLocation,
		Logger:      logger,
		Repository:  articleRepository,
		Publisher:   producers,
	})

	article.NewArticleHandler(logger, vld, router, articleUsecase)

	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	srv := server.NewServer(logger, handler, cfg.Application.Port)
	srv.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	<-sigterm

	srv.Close()
	mca.Disconnect(context.Background())

}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
