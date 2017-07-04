package service

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/NagByte/Palette/db"
	"gitlab.com/NagByte/Palette/service/auth"
	"gitlab.com/NagByte/Palette/service/checkVersion"
	"gitlab.com/NagByte/Palette/service/fileServer"
	"gitlab.com/NagByte/Palette/service/prof"
	"gitlab.com/NagByte/Palette/service/smsVerification"
)

type service interface {
	http.Handler
	URI() string
}

func StartServing() {

	// serviecs:
	var (
		neo   = db.Neo
		mongo = db.Mongo

		checkVer   = checkVersion.New(mongo)
		smsVerific = smsVerification.New(mongo)
		authen     = auth.New(smsVerific, neo)
		fileServ   = fileServer.New(mongo, authen)
		profile    = prof.New(authen, fileServ, neo)
	)
	services := []service{checkVer, smsVerific, authen, profile, fileServ}

	var handler http.Handler = handleServices(services)

	log.Fatalln(
		http.ListenAndServe(":2128", handler),
	)
}

func handleServices(services []service) *mux.Router {
	router := mux.NewRouter()
	suffix := "/{_dummy:.*}"

	for _, s := range services {
		router.Handle(s.URI()+suffix, s)
	}

	return router
}
