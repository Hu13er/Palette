package smsVerification

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	vcommon "gitlab.com/NagByte/Palette/common"
	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/helper"
	"gitlab.com/NagByte/Palette/service/common"
)

type SMSVerification interface {
	SendVerification(string) error
	Verify(string, string) (string, error)
	IsVerified(string) (string, bool)
}

type smsService struct {
	db      *smsVerificationDB
	baseURI string
	handler http.Handler
}

func New(db wrapper.Database) *smsService {
	service := &smsService{}
	service.baseURI = "/sms"
	service.db = &smsVerificationDB{db}

	router := mux.NewRouter()
	router.HandleFunc(service.baseURI+"/send/", service.sendVerificationHandler).Methods("POST")
	router.HandleFunc(service.baseURI+"/verify/", service.verifyPhoneHandler).Methods("POST")
	router.HandleFunc(service.baseURI+"/isVerified/", service.isVerifiedHandler).Methods("POST")

	service.handler = router

	service.handler = common.JSONContentTypeHandler{Handler: service.handler}

	return service
}

func (ss *smsService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ss.handler.ServeHTTP(w, r)
}

func (ss *smsService) URI() string {
	return ss.baseURI
}

/*
|* Handlers:
\********************************************/

func (ss *smsService) sendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	var form sendVerificationRequest
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	switch err := ss.SendVerification(form.PhoneNumber); err {
	case nil:
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (ss *smsService) verifyPhoneHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	var form verifyPhoneRequest
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	switch token, err := ss.Verify(form.PhoneNumber, form.VerificationCode); err {
	case nil:
		jsonEncoder.Encode(verifyPhoneResponse{Token: token})
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (ss *smsService) isVerifiedHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	var form isVerifiedRequest
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	phoneNumber, ok := ss.IsVerified(form.Token)
	jsonEncoder.Encode(isVerifiedResponse{phoneNumber, ok})
}

/*
|* API:
\***************************************/

func (ss *smsService) SendVerification(phoneNumber string) error {
	log := logrus.WithField("WHERE", "[service.smsVerification.SendVerification()]")
	var (
		code  = helper.NumricCharset.RandomStr(6)
		token = helper.DefaultCharset.RandomStr(25)
	)

	if !vcommon.ConfigBool("DEBUG") {
		_, err := sendVerification(phoneNumber, code)
		if err != nil {
			return err
		}
	} else {
		log.Debugln("verification Code:", code)
	}

	err := ss.db.mergeVerificationRequest(phoneNumber, code, token)
	return err
}

func (ss *smsService) Verify(phoneNumber, verificationCode string) (token string, err error) {
	return ss.db.verifyRequest(phoneNumber, verificationCode)
}

func (ss *smsService) IsVerified(verificationToken string) (string, bool) {
	return ss.db.isVerified(verificationToken)
}
