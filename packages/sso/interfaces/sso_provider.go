package interfaces

import (
	"net/http"

	"github.com/ESSantana/boilerplate-go/internal/domain/dto"
)

type SSOProvider interface {
	GetSigninURL() (signinURL, userState string)
	GetUserData(callbackRequest *http.Request) (data dto.UserDataSSO, err error)
}
