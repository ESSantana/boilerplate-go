package providers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/application-ellas/ellas-backend/internal/domain/dto"
	cache_interfaces "github.com/application-ellas/ellas-backend/packages/cache/interfaces"
	"github.com/application-ellas/ellas-backend/packages/sso/interfaces"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type googleSSOCallbackResponse struct {
	Id              string `json:"id"`
	Name            string `json:"given_name"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"picture"`
	VerifiedEmail   bool   `json:"verified_email"`
}

type googleSSOProvider struct {
	cacheManager cache_interfaces.CacheManager
	ssoProvider  *oauth2.Config
}

func NewGoogleSSOProvider(cacheManager cache_interfaces.CacheManager, redirectURL, clientID, secretKey string) interfaces.SSOProvider {
	return &googleSSOProvider{
		cacheManager: cacheManager,
		ssoProvider: &oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     clientID,
			ClientSecret: secretKey,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"profile", "email", "openid"},
		},
	}
}

func (provider *googleSSOProvider) GetSigninURL() (signinURL, userState string) {
	id := uuid.New().String()
	url := provider.ssoProvider.AuthCodeURL(id)
	return url, id
}

func (provider *googleSSOProvider) GetUserData(callbackRequest *http.Request) (data dto.UserDataSSO, err error) {
	state := callbackRequest.FormValue("state")
	code := callbackRequest.FormValue("code")

	if val, err := provider.cacheManager.GetFlag(callbackRequest.Context(), state); err != nil || !val {
		return data, errors.New("invalid user state")
	}

	token, err := provider.ssoProvider.Exchange(callbackRequest.Context(), code)
	if err != nil {
		return data, err
	}

	response, err := http.Get(GoogleUserInfoURL + token.AccessToken)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	var responseBody googleSSOCallbackResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return data, err
	}

	if !responseBody.VerifiedEmail {
		return data, errors.New("email not verified")
	}

	return dto.UserDataSSO{
		Email:           responseBody.Email,
		Name:            responseBody.Name,
		ExternalID:      responseBody.Id,
		ProfileImageURL: responseBody.ProfileImageURL,
	}, nil
}
