package fileServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/service/common"
)

func (fs *fileServ) uploadHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)

	contentType := r.Header.Get("Content-Type")

	size := r.Header.Get("size")
	if size == "" {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	widthAndHeight := strings.Split(size, "x")

	switch fileToken, err :=
		fs.SaveFile(
			contentType,
			map[string]interface{}{
				"width":  widthAndHeight[0],
				"height": widthAndHeight[1],
			},
			r.Body); err {
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

	switch start, c, m, err := fs.ReadFile(fileToken, w); err {
	case nil:
		w.Header().Add("Content-Type", c)

		width, height := m["width"], m["height"]
		w.Header().Add("size", fmt.Sprintf("%sx%s", width, height))

		if err := start(); err != nil {
			w.WriteHeader(common.StatusInternalServerError)
			jsonEncoder.Encode(common.ResponseInternalServerError)
		}
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}
