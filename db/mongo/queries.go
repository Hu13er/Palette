package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var queryFuncSet = map[string]queryFuncWithoutDB{
	// Check Version:
	"getAllVersion": func(db *mgo.Database, params map[string]interface{}) (*mgo.Query, slicerFunc, error) {
		return db.C("versions").Find(nil).Sort("-version"), makeSlicer("version"), nil
	},
	"getAllForcedVersion": func(db *mgo.Database, params map[string]interface{}) (*mgo.Query, slicerFunc, error) {
		return db.C("versions").Find(bson.M{"forced": true}).Sort("-version"), makeSlicer("version"), nil
	},
	// Sms Verification:
	"mergeVerificationRequest": func(db *mgo.Database, params map[string]interface{}) (*mgo.Query, slicerFunc, error) {
		_, err := db.C("smsVerfication").Upsert(bson.M{
			"phoneNumber": params["phoneNumber"],
		}, bson.M{
			"$set": bson.M{
				"code":     params["code"],
				"token":    params["token"],
				"verified": false,
				// "ttl": ....,
				// TODO: Set ttl...
			}})
		return nil, nil, err
	},
	"verifyRequest": func(db *mgo.Database, params map[string]interface{}) (*mgo.Query, slicerFunc, error) {
		err := db.C("smsVerification").Update(bson.M{
			"phoneNumber": params["phoneNumber"],
			"code":        params["code"],
		}, bson.M{
			"$set": bson.M{
				"verify": true,
			}})
		if err != nil {
			return nil, nil, err
		}

		return db.C("smsVerification").Find(bson.M{
			"phoneNumber": params["phoneNumber"],
		}), makeSlicer("token"), nil
	},
	"isVerified": func(db *mgo.Database, params map[string]interface{}) (*mgo.Query, slicerFunc, error) {
		return db.C("smsVerification").Find(bson.M{"token": params["token"]}), makeSlicer("verified", "phoneNumber"), nil
	},
}

type queryFunc func(map[string]interface{}) (*mgo.Query, slicerFunc, error)
type queryFuncWithoutDB func(*mgo.Database, map[string]interface{}) (*mgo.Query, slicerFunc, error)

type slicerFunc func(map[string]interface{}) []interface{}

func makeSlicer(keys ...string) slicerFunc {
	return func(Map map[string]interface{}) []interface{} {
		outp := make([]interface{}, len(keys))
		for i, key := range keys {
			outp[i], _ = Map[key]
		}
		return outp
	}
}

func makeAllInOneSlicer() slicerFunc {
	return func(Map map[string]interface{}) []interface{} {
		return []interface{}{Map}
	}
}

func (mongo *mongoDB) GetQuery(key string) interface{} {

	var db = mongo.db

	if v, ok := queryFuncSet[key]; ok {
		return queryFunc(func(params map[string]interface{}) (*mgo.Query, slicerFunc, error) {
			return v(db, params)
		})
	}

	panic("query " + key + " does not exist.")
}
