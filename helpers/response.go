package helpers

import (
	"encoding/json"
	"net/http"
)

const (
	RESPONSE_ERROR_MESSAGE_INTERNAL_SERVER_ERROR string = "A server error has occured, please try again later"
	RESPONSE_ERROR_MESSAGE_BAD_REQUEST_INPUT_NOT_VALID string = "The request body is not a valid JSON or doesn't match the expected structure"
)

type MinimalResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	MinimalResponse
	Message string `json:"message,omitempty"`
}

func JsonResponse( writer http.ResponseWriter, statusCode int, responseObject interface{}) {
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
