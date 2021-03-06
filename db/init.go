package db

import (
	"os"

	"github.com/sirupsen/logrus"

	"gitlab.com/NagByte/Palette/common"
	"gitlab.com/NagByte/Palette/db/mongo"
	"gitlab.com/NagByte/Palette/db/neo4j"
	"gitlab.com/NagByte/Palette/db/wrapper"
)

var (
	Neo   wrapper.Database
	Mongo wrapper.Database
)

func init() {
	mongoInit()
	neoInit()
}

func neoInit() {
	log := logrus.WithField("WHERE", "db.neo4j.neoinit")

	var err error
	neoURI := common.ConfigString("NEO_URI")
	if neoURI == "" {
		log.Errorln("Variable NEO_URI not presented.")
		os.Exit(1)
	}

	log.Infof("Connecting Neo4J on %s...\n", neoURI)

	Neo, err = neo4j.New(neoURI, 25)
	if err != nil {
		log.Panicln("Can not connect to Neo4j: ", err.Error())
	}

	Neo.Init()
}

func mongoInit() {
	log := logrus.WithField("WHERE", "[db.neo4j.mongoInit()]")
	var err error

	mongoURI := common.ConfigString("MONGO_URI")
	if mongoURI == "" {
		log.Errorln("Variable MONGO_URI not presented.")
		os.Exit(1)
	}

	log.Infof("Connecting Mongo on %s...\n", mongoURI)
	Mongo, err = mongo.New(mongoURI, "palette")
	if err != nil {
		log.Panicln("Can not connect to MongoDB: ", err.Error())
	}

	Mongo.Init()

	common.AddMongoHooker(mongoURI, "palette", "log")
}
