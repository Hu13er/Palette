package auth

/*
|* API:
\************************************/

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"

	"gitlab.com/NagByte/Palette/service/common"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

func (as *authService) DeviceTokenNeededMiddleware(f handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceToken := r.Header.Get("deviceToken")
		log.Println("in the middleware.")
		defer log.Println("out of middleware.")

		if deviceToken == "" {
			w.WriteHeader(common.StatusBadRequestError)
			json.NewEncoder(w).Encode(responseDeviceTokenNeeded)
			return
		}

		if as.ensureDeviceToken(deviceToken) == false {
			w.WriteHeader(common.StatusBadRequestError)
			json.NewEncoder(w).Encode(responseWrongDeviceToken)
			return
		}

		context.Set(r, "deviceToken", deviceToken)
		defer context.Delete(r, "deviceToken")
		f(w, r)
	}
}

func (as *authService) AuthenticationNeededMiddleware(f handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceToken := r.Header.Get("deviceToken")
		if deviceToken == "" {
			w.WriteHeader(common.StatusBadRequestError)
			json.NewEncoder(w).Encode(common.ErrorJSONResponse{ErrorDescription: "deviceTokenNeeded"})
			return
		}

		username := as.WhoAmI(deviceToken)
		if username == "" {
			w.WriteHeader(common.StatusBadRequestError)
			json.NewEncoder(w).Encode(common.ErrorJSONResponse{ErrorDescription: "authenticationNeeded"})
			return
		}

		context.Set(r, "deviceToken", deviceToken)
		defer context.Delete(r, "deviceToken")
		context.Set(r, "username", username)
		defer context.Delete(r, "username")
		f(w, r)
	}
}

func GetDeviceToken(r *http.Request) string {
	return context.Get(r, "deviceToken").(string)
}

func GetUsername(r *http.Request) string {
	return context.Get(r, "username").(string)
}
