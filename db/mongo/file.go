package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

type mongoFile struct {
	*mgo.GridFile
}

func (mf *mongoFile) ID() string { return mf.GridFile.Name() }

func (mf *mongoFile) Size() int64 { return mf.GridFile.Size() }

func (mf *mongoFile) MD5() string { return mf.GridFile.MD5() }

func (mf *mongoFile) CreatedAt() time.Time { return mf.GridFile.UploadDate() }

func (mf *mongoFile) SetContentType(contentType string) { mf.GridFile.SetContentType(contentType) }

func (mf *mongoFile) ContentType() string { return mf.GridFile.ContentType() }

func (ms *mongoFile) Metadata() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := ms.GridFile.GetMeta(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ms *mongoFile) SetMeta(meta map[string]interface{}) error {
	ms.GridFile.SetMeta(meta)
	return nil
}
