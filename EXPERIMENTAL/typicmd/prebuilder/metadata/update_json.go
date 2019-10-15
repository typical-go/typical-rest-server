package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Update the metadata
func Update(name string, v interface{}) (updated bool, err error) {
	var cachedData, data []byte
	cachedData, err = ioutil.ReadFile(Path(name))
	data, _ = json.Marshal(v)
	if os.IsNotExist(err) {
		updated = true
	} else {
		updated = (bytes.Compare(data, cachedData) != 0)
	}
	if updated {
		err = ioutil.WriteFile(Path(name), data, 0777)
	}
	return
}

// Path of metadata
func Path(name string) string {
	wd, _ := os.Getwd()
	dir := wd + "/.typical-metadata"
	os.Mkdir(dir, 0777)
	return fmt.Sprintf("%s/%s.json", dir, name)
}

// CleanAll the metadata
func CleanAll() error {
	wd, _ := os.Getwd()
	dir := wd + "/.typical-metadata"
	return os.RemoveAll(dir)
}
