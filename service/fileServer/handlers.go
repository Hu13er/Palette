package fileServer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/NagByte/Palette/service/common"
)

func (fs *fileServ) uploadHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)

	contentType := r.Header.Get("Content-Type")
	switch fileToken, err := fs.SaveFile(contentType, r.Body); err {
	case nil:
		jsonEncoder.Encode(uploadResponse{Token: fileToken})
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (fs *fileServ) downloadHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	fileToken := mux.Vars(r)["fileToken"]

	switch start, c, err := fs.ReadFile(fileToken, w); err {
	case nil:
		w.Header().Add("Content-Type", c)
		if err := start(); err != nil {
			w.WriteHeader(common.StatusInternalServerError)
			jsonEncoder.Encode(common.ResponseInternalServerError)
		}
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}
