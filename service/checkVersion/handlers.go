package checkVersion

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/service/common"
)

type handler struct {
	baseURI string
	db      database
	handler http.Handler
}

func New(db wrapper.Database) *handler {
	handl := handler{}
	handl.baseURI = "/checkVersion"
	handl.db = database{db}

	router := mux.NewRouter()
	router.HandleFunc(handl.baseURI+"/minimum/", handl.getMinimumVersionHandler).Methods("GET")
	router.HandleFunc(handl.baseURI+"/latest/", handl.getLatestVersionHandler).Methods("GET")

	handl.handler = router
	handl.handler = common.JSONContentTypeHandler{Handler: handl.handler}

	return &handl
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

func (h *handler) URI() string {
	return h.baseURI
}

func (h *handler) getMinimumVersionHandler(w http.ResponseWriter, r *http.Request) {
	if version, err := h.db.getMinimumVersion(); err == nil {
		json.NewEncoder(w).Encode(versionJSONResponse{Version: version})
	} else {
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.ResponseInternalServerError)
	}
}

func (h *handler) getLatestVersionHandler(w http.ResponseWriter, r *http.Request) {
	if version, err := h.db.getLatestVersion(); err == nil {
		json.NewEncoder(w).Encode(versionJSONResponse{Version: version})
	} else {
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.ResponseInternalServerError)
	}
}
