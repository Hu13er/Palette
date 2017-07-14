package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/service/common"
	"gitlab.com/NagByte/Palette/service/smsVerification"
)

type Auth interface {
	TouchDevice(map[string]interface{}) (string, bool, error)
	Signup(string, string, string, string) error
	SignDeviceIn(string, string, string) error
	WhoAmI(string) string
	IsUniquePhoneNumber(string) bool

	DeviceTokenNeededMiddleware(handlerFunc) handlerFunc
	AuthenticationNeededMiddleware(handlerFunc) handlerFunc
}

type authService struct {
	baseURI string
	handler http.Handler

	db       wrapper.Database
	verifier smsVerification.SMSVerification
}

func New(verifier smsVerification.SMSVerification, db wrapper.Database) *authService {
	service := &authService{}
	service.db = db
	service.verifier = verifier
	service.baseURI = "/auth"

	router := mux.NewRouter()
	router.HandleFunc(service.baseURI+"/signUp/", service.DeviceTokenNeededMiddleware(service.signUpHandler)).
		Methods("POST")
	router.HandleFunc(service.baseURI+"/touchDevice/", service.touchDeviceHandler).
		Methods("POST")
	router.HandleFunc(service.baseURI+"/signIn/", service.DeviceTokenNeededMiddleware(service.signInHandler)).
		Methods("POST")
	router.HandleFunc(service.baseURI+"/signOut/", service.DeviceTokenNeededMiddleware(service.signDeviceOutHandler)).
		Methods("POST")
	router.HandleFunc(service.baseURI+"/whoAmI/", service.DeviceTokenNeededMiddleware(service.whoAmIHandler)).
		Methods("GET")

	service.handler = router
	service.handler = common.JSONContentTypeHandler{Handler: service.handler}

	return service
}

func (as *authService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	as.handler.ServeHTTP(w, r)
}

func (as *authService) URI() string {
	return as.baseURI
}
