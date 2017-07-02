package prof

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/service/auth"
	"gitlab.com/NagByte/Palette/service/common"
)

type profService struct {
	baseURI string
	handler http.Handler
	db      wrapper.Database

	auth auth.Auth
}

func New(auth auth.Auth, db wrapper.Database) *profService {
	service := &profService{}
	service.db = db
	service.auth = auth
	service.baseURI = "/prof"

	router := mux.NewRouter()
	router.HandleFunc(service.baseURI+"/getProfile/{username}/", service.auth.AuthenticationNeededMiddleware(service.getProfileHandler)).Methods("GET")
	router.HandleFunc(service.baseURI+"/updateProfile/", service.auth.AuthenticationNeededMiddleware(service.updateProfileHandler)).Methods("POST")

	router.HandleFunc(service.baseURI+"/follow/{username}/", service.auth.AuthenticationNeededMiddleware(service.followHandler)).Methods("POST")
	router.HandleFunc(service.baseURI+"/unfollow/{username}/", service.auth.AuthenticationNeededMiddleware(service.unfollowHandler)).Methods("POST")

	router.HandleFunc(service.baseURI+"/newPost/", service.auth.AuthenticationNeededMiddleware(service.postHandler)).Methods("POST")
	router.HandleFunc(service.baseURI+"/getPosts/{username}/", service.auth.AuthenticationNeededMiddleware(service.getPostsHandler)).Methods("POST")
	router.HandleFunc(service.baseURI+"/getTimeline/", service.auth.AuthenticationNeededMiddleware(service.getTimelineHandler)).Methods("POST")

	service.handler = router
	service.handler = common.JSONContentTypeHandler{Handler: service.handler}

	return service
}

func (ps *profService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps.handler.ServeHTTP(w, r)
}

func (ps *profService) URI() string {
	return ps.baseURI
}
