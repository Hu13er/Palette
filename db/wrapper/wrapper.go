package wrapper

import (
	"io"
	"time"
)

type Rows interface {
	Columns() []string
	Next() ([]interface{}, error)
	Close() error
}

type File interface {
	io.Writer
	io.Reader
	Abort()
	io.Closer

	ID() string
	Size() int64
	MD5() string
	CreatedAt() time.Time
	SetContentType(string)
	ContentType() string
	Metadata() (map[string]interface{}, error)
	SetMeta(map[string]interface{}) error
}

type Database interface {
	Init()

	Query(interface{}, map[string]interface{}) (Rows, error)
	QueryOne(interface{}, map[string]interface{}) ([]interface{}, error)
	QueryAll(interface{}, map[string]interface{}) ([][]interface{}, error)

	Exe(interface{}, map[string]interface{}) error
	ExePipeline([]interface{}, ...map[string]interface{}) error

	GetQuery(string) interface{}

	CreateFile() (File, error)
	GetFile(string) (File, error)

	Close() error
}
