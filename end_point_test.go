package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"

	"github.com/Shopify/sarama"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jeri06/article-command-service/article"
	"github.com/jeri06/article-command-service/config"
	"github.com/jeri06/article-command-service/model"
	"github.com/jeri06/article-command-service/mongodb"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_users(t *testing.T) {
	utcLocation, _ := time.LoadLocation("UTC")
	cfg := config.Load()
	r := mux.NewRouter()
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

	articleRepository := article.NewArticleRepository(logger, mdb)
	articleUsecase := article.NewArticleUsecase(article.UsecaseProperty{
		ServiceName: cfg.Application.Name,
		UTCLoc:      utcLocation,
		Logger:      logger,
		Repository:  articleRepository,
		Publisher:   producers,
	})

	//kafka
	hh := article.HTTPHandler{
		Logger:   logger,
		Validate: vld,
		Usecase:  articleUsecase,
	}

	r.HandleFunc("/command-service/v1/article", hh.Create)
	payload := model.Article{
		ID:     1234,
		Author: "test",
		Title:  "test",
		Body:   "tes",
	}

	jsonArticle, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/command-service/v1/article", bytes.NewReader(jsonArticle))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	fmt.Println(r)

	assert.Equal(t, http.StatusOK, rr.Code)

}
