package fileServer

import (
	"fmt"
	"io"
)

type FileServer interface {
	SaveFile(contentType string, metadata map[string]interface{}, reader io.Reader) (id string, err error)
	ReadFile(id string, writer io.Writer) (starter func() error, contentType string, meta map[string]interface{}, err error)
	SmallDownloadURL(fileToken string) string
	LargeDownloadURL(fileToken string) string
}

func (fs *fileServ) SaveFile(contentType string, meta map[string]interface{}, r io.Reader) (id string, err error) {
	f, err := fs.db.CreateFile()

	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		f.Abort()
		return "", err
	}

	f.SetContentType(contentType)
	f.SetMeta(meta)

	return f.ID(), nil
}

func (fs *fileServ) ReadFile(id string, w io.Writer) (starter func() error, contentType string, meta map[string]interface{}, err error) {
	f, err := fs.db.GetFile(id)
	if err != nil {
		return nil, "", nil, err
	}

	contentType = f.ContentType()
	meta, err = f.Metadata()
	if err != nil {
		return nil, "", nil, err
	}

	return func() error {
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			return err
		}

		return nil
	}, contentType, meta, nil
}

func (fs *fileServ) SmallDownloadURL(fileToken string) string {
	if fileToken == "" {
		return ""
	}
	return fmt.Sprintf("%s/download/small/%s/", fs.baseURI, fileToken)
}

func (fs *fileServ) LargeDownloadURL(fileToken string) string {
	if fileToken == "" {
		return ""
	}
	return fmt.Sprintf("%s/download/large/%s/", fs.baseURI, fileToken)
}
