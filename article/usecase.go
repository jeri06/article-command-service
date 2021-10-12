package article

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jeri06/article-command-service/entity"
	"github.com/jeri06/article-command-service/exception"
	"github.com/jeri06/article-command-service/model"
	"github.com/jeri06/article-command-service/response"
	"github.com/sirupsen/logrus"
)

const (
	RFC3339MillisWithTripleZeroFractionSecond string = "2006-01-02T15:04:05.000Z"
	createdArticleTopic                       string = "created-article-topic"
)

type Usecase interface {
	Save(ctx context.Context, payload model.Article) (resp response.Response)
}

type usecase struct {
	serviceName string
	utcLoc      *time.Location
	logger      *logrus.Logger
	repository  Repository
	pubsub      sarama.SyncProducer
}

func NewArticleUsecase(property UsecaseProperty) Usecase {
	return &usecase{
		serviceName: property.ServiceName,
		utcLoc:      property.UTCLoc,
		logger:      property.Logger,
		repository:  property.Repository,
		pubsub:      property.Publisher,
	}
}

func (u usecase) Save(ctx context.Context, payload model.Article) (resp response.Response) {

	article := entity.Article{
		ID:      payload.ID,
		Author:  payload.Author,
		Title:   payload.Title,
		Body:    payload.Body,
		Created: time.Now().In(u.utcLoc).Format(RFC3339MillisWithTripleZeroFractionSecond),
	}
	fmt.Println(article)

	if err := u.repository.Save(ctx, article); err != nil {
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}
	articleBuff, _ := json.Marshal(article)
	message := &sarama.ProducerMessage{
		Topic: createdArticleTopic,
		Value: sarama.StringEncoder(articleBuff),
	}
	u.pubsub.SendMessage(message)
	return response.NewSuccessResponse(nil, response.StatOK, "")

}
