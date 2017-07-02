package neo4j

import (
	"os"
	"time"

	"github.com/Hu13er/logger"
)

func (DB *neo4jDB) Init() {
	DB.ensureIndexes()
	DB.ensureConstraint()
	DB.ttl()
}

func (DB *neo4jDB) ensureIndexes() {
}

func (DB *neo4jDB) ensureConstraint() {
	log := logger.WithHeaderln("[db.neo4j.ensureConstraint()]")

	err := DB.ExePipeline([]interface{}{
		"CREATE CONSTRAINT ON (user:User) ASSERT user.username IS UNIQUE",
		"CREATE CONSTRAINT ON (user:User) ASSERT user.phoneNumber IS UNIQUE",
		"CREATE CONSTRAINT ON (user:User) ASSERT user.token IS UNIQUE",
		"CREATE CONSTRAINT ON (wuser:WaitingUser) ASSERT wuser.phoneNumber IS UNIQUE",
	}, nil, nil, nil, nil)

	if err != nil {
		log.Panicln("Can not create constraint: ", err.Error())
		os.Exit(1)
	}
}

func (DB *neo4jDB) ttl() {
	go func() {
		cypher := `
		MATCH (node:TTL)
		WHERE node.ttl < timestamp()
		DETACH DELETE node;
		`
		ticker := time.NewTicker(time.Second * 60)
		for range ticker.C {
			log := logger.WithHeaderln("[db.neo4j.ttl()]")
			if err := DB.Exe(cypher, nil); err != nil {
				log.Errorln("Can not clear expired nodes: ", err.Error())
			}
		}
	}()
}
