package helpers

import (
	"encoding/json"
	"net/http"
)

type GeneralResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func JsonResponse(writer http.ResponseWriter, statusCode int, responseObject interface{}) {
	if responseObject!=nil {
		(writer).Header().Set("Content-Type", "application/json; charset=utf-8")
		(writer).WriteHeader(statusCode)
		b, _ := json.Marshal(responseObject)
		(writer).Write(b)
	} else {
		(writer).Header().Set("Content-Type", "text/plain; charset=utf-8")
		(writer).WriteHeader(statusCode)
	}

}
