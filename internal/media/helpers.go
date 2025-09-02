package media

import (
	"io"
	"net/http"
)

func GetMimeType(r io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	mimeType := http.DetectContentType(buf[:n])
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return mimeType, nil
}
