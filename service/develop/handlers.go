package develop

import (
	"net/http"
)

func (ds *developService) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}
