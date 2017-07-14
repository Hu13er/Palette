package auth

/*
|* Handlers:
\***************************************/

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/NagByte/Palette/service/common"
)

func (as *authService) touchDeviceHandler(w http.ResponseWriter, r *http.Request) {
	form := deviceProps{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		json.NewEncoder(w).Encode(common.ResponseBadRequest)
		return
	}

	validatedForm, err := form.Validate()
	if err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		json.NewEncoder(w).Encode(common.ResponseBadRequest)
		return
	}

	if deviceToken, signedIn, err := as.TouchDevice(validatedForm); err == nil {
		json.NewEncoder(w).Encode(TouchDeviceResponse{
			DeviceToken: deviceToken,
			SignedIn:    signedIn,
		})
	} else {
		w.WriteHeader(common.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.ResponseInternalServerError)
	}
}

func (as *authService) signUpHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	deviceToken := GetDeviceToken(r)

	var form signUpRequest
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
	}

	switch err := as.Signup(deviceToken, form.Username, form.Password, form.VerificationToken); err {
	case nil:
	case WrongDeviceTokenError:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(wrongDeviceTokenResponse)
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (as *authService) signInHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	deviceToken := GetDeviceToken(r)

	var form signInRequest
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
	}

	switch err := as.SignDeviceIn(deviceToken, form.Username, form.Password); err {
	case nil:
	case WrongUsernameOrPass:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(wrongUsernameOrPasswordResponse)
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (as *authService) signDeviceOutHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)

	deviceToken := GetDeviceToken(r)

	switch err := as.SignDeviceOut(deviceToken); err {
	case nil:
	case WrongDeviceTokenError:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(wrongDeviceTokenResponse)
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (as *authService) whoAmIHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)

	deviceToken := GetDeviceToken(r)

	if result := as.WhoAmI(deviceToken); result != "" {
		jsonEncoder.Encode(WhoAmIResponse{result})
	} else {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(wrongDeviceTokenResponse)
	}
}
