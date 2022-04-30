package actions

import (
	"log"
	"net/http"
)

type AvailableAction struct{}

func (aa *AvailableAction) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println(response.Write([]byte("Available Action Processed")))
}
