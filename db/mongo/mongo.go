package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/NagByte/Palette/db/wrapper"
	"gitlab.com/NagByte/Palette/helper"
)

type mongoDB struct {
	session *mgo.Session
	db      *mgo.Database
	fs      *mgo.GridFS
}

func New(uri, dbName string) (wrapper.Database, error) {
	mongo := &mongoDB{}
	var err error

	mongo.session, err = mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	mongo.db = mongo.session.DB(dbName)
	mongo.fs = mongo.db.GridFS("fs")

	return mongo, nil
}

func (db *mongoDB) QueryOne(query interface{}, params map[string]interface{}) ([]interface{}, error) {
	qf, ok := query.(queryFunc)
	if !ok {
		panic("query is not queryFunc")
	}

	q, slicer, err := qf(params)
	if err != nil {
		return nil, err
	}

	result := bson.M{}
	err = q.One(&result)

	return slicer(result), err
}

func (db *mongoDB) QueryAll(query interface{}, params map[string]interface{}) ([][]interface{}, error) {
	qf, ok := query.(queryFunc)
	if !ok {
		panic("query is not queryFunc")
	}

	q, slicer, err := qf(params)
	if err != nil {
		return nil, err
	}

	result := []bson.M{}
	q.All(&result)

	outp := make([][]interface{}, len(result))
	for i, v := range result {
		outp[i] = slicer(v)
	}

	return outp, err
}

func (db *mongoDB) Query(query interface{}, params map[string]interface{}) (wrapper.Rows, error) {
	// TODO
	panic("not implimented")
}

func (db *mongoDB) Exe(query interface{}, params map[string]interface{}) error {
	qf, ok := query.(queryFunc)
	if !ok {
		panic("query is not queryFunc")
	}

	_, _, err := qf(params)
	return err
}

func (db *mongoDB) ExePipeline([]interface{}, ...map[string]interface{}) error {
	// TODO
	panic("not implimented")
}

func (db *mongoDB) CreateFile() (wrapper.File, error) {
	id := helper.DefaultCharset.RandomStr(30)
	gf, err := db.fs.Create(id)
	if err != nil {
		return nil, err
	}

	return &mongoFile{gf}, nil
}

func (db *mongoDB) GetFile(id string) (wrapper.File, error) {
	gf, err := db.fs.Open(id)
	if err != nil {
		return nil, err
	}

	return &mongoFile{gf}, nil
}

func (db *mongoDB) Close() error {
	db.session.Close()
	return nil
}
