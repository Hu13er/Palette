package fileServer

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/service/auth"
	"gitlab.com/NagByte/Palette/service/common"
)

type fileServ struct {
	baseURI string
	handler http.Handler

	db     wrapper.Database
	authen auth.Auth
}

func New(db wrapper.Database, auth auth.Auth) *fileServ {
	fs := &fileServ{}
	fs.db = db
	fs.authen = auth
	fs.baseURI = "/storage"

	router := mux.NewRouter()
	router.HandleFunc(fs.baseURI+"/upload/", common.AddJSONContentType(fs.uploadHandler)).
		Methods("POST")
	router.HandleFunc(fs.baseURI+"/download/small/{fileToken}/", fs.downloadHandler).
		Methods("GET")
	router.HandleFunc(fs.baseURI+"/download/large/{fileToken}/", fs.downloadHandler).
		Methods("GET")

	fs.handler = router

	return fs
}

func (fs *fileServ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.handler.ServeHTTP(w, r)
}

func (fs *fileServ) URI() string {
	return fs.baseURI
}
