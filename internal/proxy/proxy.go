package proxy

import (
	"fmt"
	"lproxy/internal/auth"
	"lproxy/pkg/log"
	"net/http"
	"strings"
)

type Proxy struct {
	port   string
	auth   *auth.Auth
	points map[string]*Point
	errCh  chan log.Message
}

func InitProxy(config *Config, auth *auth.Auth, errCh chan log.Message) *Proxy {
	proxy := &Proxy{
		port:   config.Port,
		auth:   auth,
		points: make(map[string]*Point, 0),
		errCh: errCh,
	}
	proxy.initPoints(config.Points)
	return proxy
}

func (p *Proxy) initPoints(pointConfigurations []*PointConfig) {
	for _, config := range pointConfigurations {
		if _, ok := p.points[config.Code]; ok {
			//glog.Exitf("Cannot add point with name: %s. Point with same code is already exists.", config.Code)
		}
		pointProviders := p.auth.GetProvidersFromList(config.Providers)
		if len(pointProviders) != len(config.Providers) {
			//glog.Exitf("Cannot identify all point's authorization providers with code: %s", config.Code)
		}
		p.points[config.Code] = PointInit(config, pointProviders)
	}
}

func (p *Proxy) Stop() {}

func (p *Proxy) findPoint(request *http.Request) (*Point, error) {
	for _, point := range p.points {
		if point.isMatch(request) {
			return point, nil
		}
	}
	return nil, &Error{
		Message: fmt.Sprintf("Unknown router for request host: %s path: %s", request.Host, request.URL.Path)}
}

func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	request.Host = strings.ReplaceAll(request.Host, ":"+p.port, "")
	point, err := p.findPoint(request)
	if err != nil {
		//glog.Errorf("Unknown point. Error: %s", err.Error())
		writer.WriteHeader(404)
		_, wErr := writer.Write([]byte("Unknown proxy point"))
		if wErr != nil {
			//glog.Error(wErr)
		}
		return
	}
	executionErr := point.execute(request, writer)

	if executionErr != nil {
		//glog.Error("Cannot execute data for point: %s error: %s", point.code, executionErr.Error())
		writer.WriteHeader(500)
		_, wErr := writer.Write([]byte("Error for execute proxy"))
		if wErr != nil {
			//glog.Error(wErr)
		}
		return
	}
}
