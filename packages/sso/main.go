package sso

import (
	"errors"

	cache_interfaces "github.com/application-ellas/ella-backend/packages/cache/interfaces"
	"github.com/application-ellas/ella-backend/packages/sso/interfaces"
	"github.com/application-ellas/ella-backend/packages/sso/providers"
)

type GoogleProvider struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
}

type SSOManager struct {
	cacheManager cache_interfaces.CacheManager
	google       GoogleProvider
}

func NewSSOManager(cacheManager cache_interfaces.CacheManager, google GoogleProvider) interfaces.SSOManager {
	return &SSOManager{
		cacheManager: cacheManager,
		google:       google,
	}
}

func (s *SSOManager) GetProvider(provider string) (ssoProvider interfaces.SSOProvider, err error) {
	switch provider {
	case "google":
		return providers.NewGoogleSSOProvider(s.cacheManager, s.google.RedirectURL, s.google.ClientID, s.google.ClientSecret), nil
	default:
		return nil, errors.New("invalid provider")
	}
}
