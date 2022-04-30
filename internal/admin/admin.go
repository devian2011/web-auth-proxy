package admin

import (
	"github.com/gorilla/mux"
	"lproxy/internal/admin/actions"
	"lproxy/internal/auth"
	"net/http"
)

type Area struct {
}

func InitAdminArea(config *Config, auth *auth.Auth) *Area {
	return &Area{}
}

func (a *Area) GetHandler() http.Handler {
	router := &mux.Router{}
	router.Handle("/", &actions.AvailableAction{})

	return router
}

func (a *Area) Stop() {

}
