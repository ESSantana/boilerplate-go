package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/application-ellas/ellas-backend/internal/domain/dto"
)

func CreateResponse(response *http.ResponseWriter, statusCode int, err error, data ...any) {
	var body dto.HttpResponse
	if err != nil {
		body = dto.HttpResponse{
			Error:   true,
			Message: err.Error(),
		}
	} else {
		body = dto.HttpResponse{
			Error: false,
			Data:  data[0],
		}
	}

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
		CreateResponse(&response, http.StatusBadRequest, err)
		return bodyRequest
	}

	err = json.Unmarshal(body, &bodyRequest)
	if err != nil {
		CreateResponse(&response, http.StatusBadRequest, err)
		return bodyRequest
	}
	return bodyRequest
}
