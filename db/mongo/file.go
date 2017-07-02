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
