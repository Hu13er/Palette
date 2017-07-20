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
	// Upload images
	router.HandleFunc(
		fs.baseURI+"/",
		common.AddJSONContentType(fs.uploadHandler),
	).Methods("PUT")

	// Download small size
	router.HandleFunc(
		fs.baseURI+"/{fileToken}/small/",
		fs.downloadHandler,
	).Methods("GET")

	// Download Large size
	router.HandleFunc(
		fs.baseURI+"/{fileToken}/large/",
		fs.downloadHandler,
	).Methods("GET")

	fs.handler = router

	return fs
}

func (fs *fileServ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.handler.ServeHTTP(w, r)
}

func (fs *fileServ) URI() string {
	return fs.baseURI
}
