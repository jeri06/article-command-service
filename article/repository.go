package article

import (
	"context"

	"github.com/jeri06/article-command-service/entity"
	"github.com/jeri06/article-command-service/exception"
	"github.com/jeri06/article-command-service/mongodb"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	Save(ctx context.Context, article entity.Article) (err error)
}

type repository struct {
	logger     *logrus.Logger
	collection mongodb.Collection
}

func NewArticleRepository(logger *logrus.Logger, db mongodb.Database) Repository {
	var collectionName = "article"
	var collection mongodb.Collection

	collection = db.Collection(collectionName)

	return &repository{logger, collection}
}

func (r repository) Save(ctx context.Context, article entity.Article) (err error) {
	_, err = r.collection.InsertOne(ctx, article)
	if err != nil {
		r.logger.Error(err)
		err = exception.ErrInternalServer
		return
	}
	return
}
