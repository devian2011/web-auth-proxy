package admin

import (
	"github.com/gorilla/mux"
	"lproxy/internal/admin/actions"
	"lproxy/internal/auth"
	"lproxy/pkg/log"
	"net/http"
)

type Area struct {
	config *Config
	auth   *auth.Auth
	errCh  chan log.Message
}

func InitAdminArea(config *Config, auth *auth.Auth, errCh chan log.Message) *Area {
	return &Area{
		config: config,
		auth:   auth,
		errCh:  errCh,
	}
}

func (a *Area) GetHandler() http.Handler {
	router := &mux.Router{}
	router.Handle("/", &actions.AvailableAction{})

	return router
}

func (a *Area) Stop() {

}
