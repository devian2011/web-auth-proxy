package auth

import "lproxy/internal/auth/providers"

type Config struct {
	Providers []*providers.Config `json:"providers"`
}

