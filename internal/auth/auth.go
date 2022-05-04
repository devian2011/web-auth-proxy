package auth

import (
	"fmt"
	"lproxy/internal/auth/providers"
	"lproxy/internal/auth/providers/db"
	"lproxy/pkg/log"
)

type authError struct {
	Message string
}

func (ae *authError) Error() string {
	return ae.Message
}

type Auth struct {
	providers map[string]providers.Provider
	errCh     chan log.Message
}

func NewAuth(config *Config, errCh chan log.Message) *Auth {
	auth := &Auth{
		providers: make(map[string]providers.Provider, 0),
		errCh:     errCh,
	}
	for _, v := range config.Providers {
		provider, err := factory(v)
		if err != nil {
			errCh <- log.NewCriticalMessage(
				fmt.Sprintf("Cannot create authorization provider: %s", v.Code), 32, "auth.go", err)
		} else {
			if _, ok := auth.providers[v.Code]; ok {
				errCh <- log.NewCriticalMessage(
					fmt.Sprintf("Authentication provider with code %s is already exists", v.Code), 36, "auth.go", err)
			} else {
				auth.providers[v.Code] = provider
			}
		}
	}

	return auth
}

func (a *Auth) GetProvider(code string) (providers.Provider, error) {
	if _, ok := a.providers[code]; ok {
		return a.providers[code], nil
	}
	return nil, &authError{Message: fmt.Sprintf("Unknown provider with code: %s", code)}
}

func (a *Auth) GetProvidersFromList(codes []string) map[string]providers.Provider {
	result := make(map[string]providers.Provider, 0)
	for _, code := range codes {
		if v, ok := a.providers[code]; ok {
			result[code] = v
		}
	}

	return result
}

func (a *Auth) GetProviderList() []string {
	var result []string
	for k, _ := range a.providers {
		result = append(result, k)
	}
	return result
}

func (a *Auth) Stop() {
	for _, v := range a.providers {
		if e := v.Shutdown(); e != nil {
			a.errCh <- log.NewErrorMessage(
				fmt.Sprintf("Error on shutdown auth provider: %s", v.GetCode()), 76, "auth.go", e)
		}

	}
}

func factory(config *providers.Config) (providers.Provider, error) {
	switch config.Type {
	case db.ProviderType:
		if config.IsActive {
			return db.NewDbProvider(config)
		} else {
			return nil, nil
		}
	default:
		return nil, &authError{Message: fmt.Sprintf("Unknown provider with type: %s", config.Type)}
	}
}
