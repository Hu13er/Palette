package prof

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/service/auth"
	"gitlab.com/NagByte/Palette/service/common"
	"gitlab.com/NagByte/Palette/service/fileServer"
)

type profService struct {
	baseURI string
	handler http.Handler
	db      wrapper.Database

	auth auth.Auth
	fs   fileServer.FileServer
}

func New(authen auth.Auth, fs fileServer.FileServer, db wrapper.Database) *profService {
	service := &profService{}
	service.db = db
	service.auth = authen
	service.fs = fs
	service.baseURI = "/prof"

	router := mux.NewRouter()

	// Gets this profile
	router.HandleFunc(
		service.baseURI+"/profile/",
		service.auth.AuthenticationNeededMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				mux.Vars(r)["username"] = auth.GetUsername(r)
				service.getProfileHandler(w, r)
			}),
	).Methods("GET")

	// Update this profile
	router.HandleFunc(
		service.baseURI+"/profile/",
		service.auth.AuthenticationNeededMiddleware(service.updateProfileHandler),
	).Methods("POST")

	// Gets a usernames's profile
	router.HandleFunc(
		service.baseURI+"/profile/{username}/",
		service.auth.AuthenticationNeededMiddleware(service.getProfileHandler),
	).Methods("GET")

	// Gets timeline
	router.HandleFunc(
		service.baseURI+"/timeline/",
		service.auth.AuthenticationNeededMiddleware(service.getTimelineHandler),
	).Methods("GET")

	// Gets this user's posts
	router.HandleFunc(
		service.baseURI+"/posts/",
		service.auth.AuthenticationNeededMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				mux.Vars(r)["username"] = auth.GetUsername(r)
				service.getPostsHandler(w, r)
			}),
	).Methods("GET")

	// Gets username's posts
	router.HandleFunc(
		service.baseURI+"/posts/{username}/",
		service.auth.AuthenticationNeededMiddleware(service.getPostsHandler),
	).Methods("GET")

	// Adds new post
	router.HandleFunc(
		service.baseURI+"/posts/",
		service.auth.AuthenticationNeededMiddleware(service.postHandler),
	).Methods("PUT")

	// Like and Dislike a post
	router.HandleFunc(
		service.baseURI+"/post/{artID}/likes/",
		service.auth.AuthenticationNeededMiddleware(service.likePostHandler),
	).Methods("PUT", "DELETE")

	// Follow and unfollow someone
	router.HandleFunc(
		service.baseURI+"/follow/{username}/",
		service.auth.AuthenticationNeededMiddleware(service.followHandler),
	).Methods("PUT", "DELETE")

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
