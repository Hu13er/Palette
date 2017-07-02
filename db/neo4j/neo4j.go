package neo4j

import (
	"io"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"

	"gitlab.com/NagByte/Palette/db/wrapper"
)

var (
	ErrNotFound = io.EOF
)

type neo4jRows struct {
	bolt.Rows
	bolt.Conn
}

func (nr *neo4jRows) Columns() []string {
	return nr.Rows.Columns()
}

func (nr *neo4jRows) Next() (result []interface{}, err error) {
	result, _, err = nr.Rows.NextNeo()
	result = neoToStdEntity(result).([]interface{})
	return
}

func (nr *neo4jRows) Close() error {
	if err := nr.Rows.Close(); err != nil {
		return err
	}

	if err := nr.Conn.Close(); err != nil {
		return err
	}

	return nil
}

type neo4jDB struct {
	connPool bolt.DriverPool
}

func New(neoURI string, connCount int) (wrapper.Database, error) {
	pool, err := bolt.NewDriverPool(neoURI, connCount)
	if err != nil {
		return nil, err
	}

	return &neo4jDB{connPool: pool}, nil
}

func (ndb *neo4jDB) Query(cypher interface{}, params map[string]interface{}) (wrapper.Rows, error) {
	params = stdToNeoEntity(params)

	conn, err := ndb.connPool.OpenPool()
	if err != nil {
		return nil, err
	}

	rows, err := conn.QueryNeo(cypher.(string), params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return &neo4jRows{Rows: rows, Conn: conn}, nil
}

func (ndb *neo4jDB) QueryOne(cypher interface{}, params map[string]interface{}) ([]interface{}, error) {
	params = stdToNeoEntity(params)

	conn, err := ndb.connPool.OpenPool()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryNeo(cypher.(string), params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, _, err := rows.NextNeo()
	return neoToStdEntity(result).([]interface{}), err
}

func (ndb *neo4jDB) QueryAll(cypher interface{}, params map[string]interface{}) ([][]interface{}, error) {
	params = stdToNeoEntity(params)

	conn, err := ndb.connPool.OpenPool()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, _, _, err := conn.QueryNeoAll(cypher.(string), params)
	return neoToStdEntity(rows).([][]interface{}), err
}

func (ndb *neo4jDB) Exe(cypher interface{}, params map[string]interface{}) error {
	params = stdToNeoEntity(params)

	conn, err := ndb.connPool.OpenPool()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecNeo(cypher.(string), params)
	return err
}

func (ndb *neo4jDB) ExePipeline(cyphers []interface{}, params ...map[string]interface{}) error {
	conn, err := ndb.connPool.OpenPool()
	if err != nil {
		return err
	}
	defer conn.Close()

	strCyphers := make([]string, len(cyphers))
	for i, v := range cyphers {
		strCyphers[i] = v.(string)
	}

	_, err = conn.ExecPipeline(strCyphers, params...)
	return err
}

func (ndb *neo4jDB) CreateFile() (wrapper.File, error) {
	panic("feature not supported")
}

func (ndb *neo4jDB) GetFile(string) (wrapper.File, error) {
	panic("feature not supported")
}

func (ndb *neo4jDB) Close() error {
	return nil
}

func stdToNeoEntity(raw map[string]interface{}) map[string]interface{} {
	for k, v := range raw {
		switch v.(type) {
		case []string:
			conv := raw[k].([]string)
			outp := make([]interface{}, len(conv))
			for i, v := range conv {
				outp[i] = v
			}
			raw[k] = outp
		default:
		}
	}
	return raw
}

func neoToStdEntity(raw interface{}) interface{} {
	switch raw.(type) {
	case graph.Node:
		return raw.(graph.Node).Properties
	case []interface{}:
		conv := raw.([]interface{})
		outp := make([]interface{}, len(conv))
		for i, v := range conv {
			outp[i] = neoToStdEntity(v)
		}
		return outp
	case [][]interface{}:
		conv := raw.([][]interface{})
		outp := make([][]interface{}, len(conv))
		for i, v := range conv {
			outp[i] = neoToStdEntity(v).([]interface{})
		}
		return outp
	default:
		// TODO: use reflection package to Standardize Database output here.
		return raw
	}
}
