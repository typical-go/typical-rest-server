package typimeta

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// MetadataHandler to handling the metadata
type MetadataHandler struct {
}

func (h *MetadataHandler) Read(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var n int64 = bytes.MinRead
	if fi, err := f.Stat(); err == nil {
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return ioutil.ReadAll(f)
}

func (h *MetadataHandler) Write(filename string, data []byte) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
