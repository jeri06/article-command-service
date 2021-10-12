package article

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type UsecaseProperty struct {
	ServiceName string
	UTCLoc      *time.Location
	Logger      *logrus.Logger
	Publisher   sarama.SyncProducer
	Repository  Repository
}
