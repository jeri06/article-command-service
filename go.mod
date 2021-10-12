module github.com/jeri06/article-command-service

go 1.15

require (
	github.com/Shopify/sarama v1.29.0
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gorilla/mux v1.8.0
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/rs/cors v1.8.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.7.3
	golang.org/x/net v0.0.0-20210427231257-85d9c07bbe3a
)

// replace github.com/jeri06/article-command-service => ../
