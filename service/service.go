package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"gitlab.com/NagByte/Palette/db"
	"gitlab.com/NagByte/Palette/service/auth"
	"gitlab.com/NagByte/Palette/service/checkVersion"
	"gitlab.com/NagByte/Palette/service/develop"
	"gitlab.com/NagByte/Palette/service/fileServer"
	"gitlab.com/NagByte/Palette/service/prof"
	"gitlab.com/NagByte/Palette/service/smsVerification"
)

type service interface {
	http.Handler
	URI() string
}

func StartServing() {
	log := logrus.WithField("WHERE", "[service.service.StartServing()]")
	// serviecs:
	var (
		neo   = db.Neo
		mongo = db.Mongo
	)

	checkVer := checkVersion.New(mongo)

	smsVerific := smsVerification.New(mongo)
	authen := auth.New(smsVerific, neo)
	smsVerific.SetUniquer(authen)

	fileServ := fileServer.New(mongo, authen)
	profile := prof.New(authen, fileServ, neo)
	dev := develop.New()

	services := []service{checkVer, smsVerific, authen, profile, fileServ, dev}

	var handler http.Handler = handleServices(services)
	log.Infoln("Start Serving")

	writer := logrus.StandardLogger().WriterLevel(logrus.DebugLevel)
	handler = requestURILoggerMiddleware(writer, handler)

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
