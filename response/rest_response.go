package response

import (
	"encoding/json"
	"net/http"
)

// REST is a collection of behavior of REST.
type REST interface {
	JSON(w http.ResponseWriter)
}

type restObject struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"` // will not be appeared if not set.
	// can add more
}

// JSON will response as json serialization.
func JSON(w http.ResponseWriter, resp Response) {
	var success bool
	if resp.Error() == nil {
		success = true
	}
	ro := restObject{
		Success: success,
		Data:    resp.Data(),
		Message: resp.Message(),
		Status:  resp.Status(),
		Code:    resp.HTTPStatusCode(),
		Meta:    resp.Meta(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ro.Code)
	json.NewEncoder(w).Encode(ro)
}
