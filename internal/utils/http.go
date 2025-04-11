package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/application-ellas/ella-backend/internal/domain/dto"
	"github.com/application-ellas/ella-backend/internal/domain/errors"
)

func CreateResponse(response *http.ResponseWriter, statusCode int, responseErr error, data ...any) {
	var body dto.HttpResponse
	if responseErr != nil {
		body = dto.HttpResponse{
			Error:   true,
			Message: responseErr.Error(),
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

	if _, ok := responseErr.(*errors.ValidationError); ok {
		statusCode = http.StatusUnprocessableEntity
	}

	if _, ok := responseErr.(*errors.NotFoundError); ok {
		statusCode = http.StatusNotFound
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
