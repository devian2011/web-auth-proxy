package proxy

import (
	"bytes"
	"github.com/google/uuid"
	"io/ioutil"
	"lproxy/internal/auth/providers"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Point struct {
	code          string
	authProviders map[string]providers.Provider
	config        *PointConfig
}

func PointInit(config *PointConfig, providers map[string]providers.Provider) *Point {
	return &Point{
		code:          config.Code,
		authProviders: providers,
		config:        config,
	}
}

func (p *Point) isMatch(r *http.Request) bool {
	matchHost, errHost := regexp.Match(p.config.Match.Host, []byte(r.Host))
	if errHost != nil {
		//glog.Errorf(
		//	"Regexp error for match point: %s regexp: %s with request: %s Error: %s",
		//	p.code, p.config.Match.Host, r.Host, errHost.Error())
	}
	if matchHost {
		matchPath, errPath := regexp.Match(p.config.Match.Path, []byte(r.URL.Path))
		if errPath != nil {
			//glog.Errorf(
			//	"Regexp error for match point: %s path: %s with request: %s Error: %s",
			//	p.code, p.config.Match.Path, r.URL.Path, errPath.Error())
		}
		return matchPath
	}

	return false
}

func (p *Point) execute(r *http.Request, w http.ResponseWriter) error {
	subRequest, subRequestBuildErr := p.buildSubRequest(r)
	if subRequestBuildErr != nil {
		return subRequestBuildErr
	}
	client := http.Client{}
	response, responseErr := client.Do(subRequest)
	if responseErr != nil {
		return responseErr
	}
	writeErr := write(&w, response)
	if writeErr != nil {
		return responseErr
	}

	return nil
}

func (p *Point) buildSubRequest(r *http.Request) (*http.Request, error) {
	to, destinationParseError := url.Parse(p.config.Destination)
	if destinationParseError != nil {
		return nil, destinationParseError
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	u := &url.URL{
		Scheme:     to.Scheme,
		Opaque:     r.URL.Opaque,
		User:       r.URL.User,
		Host:       to.Host,
		Path:       r.URL.Path,
		RawPath:    r.URL.RawPath,
		ForceQuery: r.URL.ForceQuery,
		RawQuery:   r.URL.RawQuery,
		Fragment:   r.URL.Fragment,
	}

	request, _ := http.NewRequest(r.Method, u.String(), ioutil.NopCloser(bytes.NewReader(body)))

	request.Header = make(http.Header)
	for h, val := range r.Header {
		request.Header[h] = val
	}
	r.Header.Add("PR-REQUEST-ID", uuid.New().String())

	return request, nil
}

func write(w *http.ResponseWriter, r *http.Response) error {
	body, readerError := ioutil.ReadAll(r.Body)
	if readerError != nil {
		return readerError
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	for k, v := range r.Header {
		(*w).Header().Add(k, strings.Join(v, ";"))
	}
	(*w).WriteHeader(200)
	_, writeError := (*w).Write(body)
	if writeError != nil {
		return writeError
	}

	return nil
}
