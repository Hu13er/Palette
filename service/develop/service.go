package develop

import (
	"net/http"

	"github.com/gorilla/mux"
)

type developService struct {
	handler http.Handler
	baseURI string
}

func New() *developService {
	service := &developService{}
	service.baseURI = "/develop"

	router := mux.NewRouter()
	router.HandleFunc(service.baseURI+"/ping/", service.pingHandler)

	service.handler = router

	return service
}

func (ds *developService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ds.handler.ServeHTTP(w, r)
}

func (ds *developService) URI() string {
	return ds.baseURI
}
