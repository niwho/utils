package utils

import (
	"compress/gzip"
	"io"
)

func Ungzip(r io.Reader) (io.ReadCloser, error) {
	rb, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return rb, nil
}
