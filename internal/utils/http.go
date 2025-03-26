package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func CreateResponse(response *http.ResponseWriter, statusCode int, body map[string]interface{}) {
	bodyResponse, err := json.Marshal(body)
	if err != nil {
		panic(errors.New("something went wrong on http utils "))
	}
	(*response).WriteHeader(statusCode)
	(*response).Write(bodyResponse)
	(*response).Header().Set("Content-Type", "application/json")
}

func ReadBody[T any](request *http.Request, response http.ResponseWriter) (output T) {
	var bodyRequest T
	body, err := io.ReadAll(request.Body)
	if err != nil {
		responseBody := map[string]interface{}{
			"message": "payload malformed",
		}
		CreateResponse(&response, http.StatusBadRequest, responseBody)
		return bodyRequest
	}

	err = json.Unmarshal(body, &bodyRequest)
	if err != nil {
		responseBody := map[string]interface{}{
			"message": "payload malformed",
		}
		CreateResponse(&response, http.StatusBadRequest, responseBody)
		return bodyRequest
	}
	return bodyRequest
}
