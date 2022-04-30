package internal

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"lproxy/internal/admin"
	"lproxy/internal/auth"
	"lproxy/internal/proxy"
	"lproxy/pkg/file"
	"net/http"
)

type Application struct {
	configuration *Configuration
	ctx           context.Context
	auth          *auth.Auth
	proxy         *proxy.Proxy
	admin         *admin.Area
}

func NewApplication(configurationFilePath string, ctx context.Context) *Application {
	configuration := parseConfiguration(configurationFilePath)
	authorization := auth.NewAuth(configuration.Auth)
	app := &Application{
		configuration: configuration,
		auth:          authorization,
		proxy:         proxy.InitProxy(configuration.Proxy, authorization),
		admin:         admin.InitAdminArea(configuration.Admin, authorization),
		ctx:           ctx,
	}

	return app
}

func parseConfiguration(configurationFilePath string) *Configuration {
	config := &Configuration{}
	configReaderErr := file.LoadStructureFromJsonFile(configurationFilePath, config)
	if configReaderErr != nil {
		glog.Exit(configReaderErr)
	}

	return config
}

func (a *Application) Run() {
	go func() {
		glog.Exit(http.ListenAndServe(
			fmt.Sprintf("%s:%s",
				a.configuration.Proxy.Host,
				a.configuration.Proxy.Port),
			a.proxy))
	}()
	go func() {
		glog.Exit(http.ListenAndServe(
			fmt.Sprintf("%s:%s",
				a.configuration.Admin.Host,
				a.configuration.Admin.Port),
			a.admin.GetHandler()))
	}()
	<- a.ctx.Done()

}

func (a *Application) stop() {
	a.proxy.Stop()
	a.admin.Stop()
	a.auth.Stop()
}
