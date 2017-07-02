package fileServer

import (
	"io"
)

func (fs *fileServ) SaveFile(contentType string, r io.Reader) (string, error) {
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
	return f.ID(), nil
}

func (fs *fileServ) ReadFile(id string, w io.Writer) (func() error, string, error) {
	f, err := fs.db.GetFile(id)
	if err != nil {
		return nil, "", err
	}

	return func() error {
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			return err
		}

		return nil
	}, f.ContentType(), nil
}
