package internal

import (
	"context"
	"fmt"
	"lproxy/internal/admin"
	"lproxy/internal/auth"
	"lproxy/internal/proxy"
	"lproxy/pkg/log"
	"net/http"
	"time"
)

type Application struct {
	configuration *Configuration
	ctx           context.Context
	auth          *auth.Auth
	proxy         *proxy.Proxy
	admin         *admin.Area
	logger        *log.Handler
}

func NewApplication(configurationFilePath string, ctx context.Context) (*Application, error) {
	configuration, err := ParseConfiguration(configurationFilePath)
	if err != nil {
		return nil, err
	}
	authorization := auth.NewAuth(configuration.Auth)
	app := &Application{
		configuration: configuration,
		auth:          authorization,
		proxy:         proxy.InitProxy(configuration.Proxy, authorization),
		admin:         admin.InitAdminArea(configuration.Admin, authorization),
		ctx:           ctx,
		logger:        initLogger(configuration),
	}

	return app, nil
}

func initLogger(configuration *Configuration) *log.Handler {
	handler := log.NewHandler()
	if configuration.Logs.File != nil {
		handler.AddLogger(
			configuration.Logs.File.Level,
			log.NewFileLogger(log.StrFormatter, configuration.Logs.File.Path))
	}
	if configuration.Logs.Console != nil {
		handler.AddLogger(
			configuration.Logs.Console.Level,
			log.NewStdoutLogger(log.StrFormatter))
	}

	return handler
}

func (a *Application) Run() {
	proxyServer := a.proxyServerStart()
	adminServer := a.adminServerStart()
	<-a.ctx.Done()
	proxyServerCtxStop, fnProxyStop := context.WithTimeout(a.ctx, 5 * time.Second)
	_ = proxyServer.Shutdown(proxyServerCtxStop)
	adminServerCtxStop, fnAdminStop := context.WithTimeout(a.ctx, time.Second)
	_ = adminServer.Shutdown(adminServerCtxStop)

	defer func() {
		fnProxyStop()
		fnAdminStop()
	}()
}

func (a *Application) proxyServerStart() *http.Server {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s",
			a.configuration.Proxy.Host,
			a.configuration.Proxy.Port),
		Handler:           a.proxy,
		TLSConfig:         nil,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	go func() {
		var err error
		if a.configuration.Proxy.Tls != nil {
			err = server.ListenAndServeTLS(a.configuration.Proxy.Tls.CertFile, a.configuration.Proxy.Tls.KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			a.logger.Handle(log.NewCriticalMessage("Cannot start proxy server", 102, "app.go", err))
		}
	}()

	return server
}

func (a *Application) adminServerStart() *http.Server {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s",
			a.configuration.Admin.Host,
			a.configuration.Admin.Port),
		Handler:           a.admin.GetHandler(),
		TLSConfig:         nil,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	go func() {
		var err error
		if a.configuration.Admin.Tls != nil {
			err = server.ListenAndServeTLS(a.configuration.Admin.Tls.CertFile, a.configuration.Admin.Tls.KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			a.logger.Handle(log.NewCriticalMessage("Cannot start admin server", 135, "app.go", err))
		}
	}()

	return server
}

func (a *Application) stop() {
	a.proxy.Stop()
	a.admin.Stop()
	a.auth.Stop()
}
