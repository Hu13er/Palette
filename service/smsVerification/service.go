package smsVerification

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	vcommon "gitlab.com/NagByte/Palette/common"
	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/helper"
	"gitlab.com/NagByte/Palette/service/common"
)

var (
	ErrPhoneNumberExists                  = errors.New("phoneNumberExists")
	ErrWrongPhoneNumberOrVerificationCode = errors.New("wrongPhoneNumberOrVerificationCode")
)

type SMSVerification interface {
	SendVerification(string, bool) error
	Verify(string, string) (string, error)
	IsVerified(string) (string, bool)

	SetUniquer(unique uniquer)
	IsUnique(phoneNumber string) bool
}

type smsService struct {
	db      *smsVerificationDB
	baseURI string
	handler http.Handler
	unique  uniquer
}

type uniquer interface {
	IsUniquePhoneNumber(phoneNumber string) bool
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

	switch err := ss.SendVerification(form.PhoneNumber, form.SignUpState); err {
	case nil:
	case ErrUnsuccessfulRequest:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseInvalidPhoneNumber)
	case ErrPhoneNumberExists:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responsePhoneNumberExists)
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
	case ErrWrongPhoneNumberOrVerificationCode:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseWrongPhoneNumberOrVerificationCode)
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

func (ss *smsService) SendVerification(phoneNumber string, signUpState bool) error {
	log := logrus.WithField("WHERE", "[service.smsVerification.SendVerification()]")

	if signUpState && !ss.IsUnique(phoneNumber) {
		return ErrPhoneNumberExists
	}

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

func (ss *smsService) SetUniquer(unique uniquer) {
	ss.unique = unique
}

func (ss *smsService) IsUnique(phoneNumber string) bool {
	return ss.unique.IsUniquePhoneNumber(phoneNumber)
}

func (ss *smsService) Verify(phoneNumber, verificationCode string) (token string, err error) {
	token, err = ss.db.verifyRequest(phoneNumber, verificationCode)
	switch err {
	case nil:
		return token, nil
	case io.EOF:
		return "", ErrWrongPhoneNumberOrVerificationCode
	default:
		return "", err
	}
}

func (ss *smsService) IsVerified(verificationToken string) (string, bool) {
	return ss.db.isVerified(verificationToken)
}
