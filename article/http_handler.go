package article

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jeri06/article-command-service/model"
	"github.com/jeri06/article-command-service/response"
	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usecase  Usecase
}

func NewArticleHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, usecase Usecase) {
	handle := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/command-service/v1/article", handle.Create).Methods(http.MethodPost)
}

func (handler HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var article model.Article

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		resp = response.NewErrorResponse(err, http.StatusUnprocessableEntity, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	if err := handler.validateRequestBody(article); err != nil {
		resp = response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}
	resp = handler.Usecase.Save(ctx, article)
	response.JSON(w, resp)
	return
}

func (handler HTTPHandler) validateRequestBody(body interface{}) (err error) {
	err = handler.Validate.Struct(body)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
