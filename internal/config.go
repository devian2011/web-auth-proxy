package internal

import (
	"lproxy/internal/admin"
	"lproxy/internal/auth"
	"lproxy/internal/proxy"
)

type Configuration struct {
	Proxy *proxy.Config `json:"proxy"`
	Auth  *auth.Config  `json:"auth"`
	Admin *admin.Config `json:"admin"`
}
