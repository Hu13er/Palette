package common

import "net/http"

type JSONContentTypeHandler struct {
	Handler http.Handler
}

func (jcth JSONContentTypeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	jcth.Handler.ServeHTTP(w, r)
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func AddJSONContentType(f handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		f(w, r)
	}
}
