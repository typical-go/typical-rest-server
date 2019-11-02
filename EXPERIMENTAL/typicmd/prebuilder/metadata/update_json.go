package metadata

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

// Update the metadata
func Update(name string, v interface{}) (updated bool, err error) {
	filename := Path(name + ".json")
	var cachedData, data []byte
	cachedData, err = ioutil.ReadFile(filename)
	if data, err = json.Marshal(v); err != nil {
		return
	}
	if os.IsNotExist(err) {
		updated = true
	} else {
		updated = (bytes.Compare(data, cachedData) != 0)
	}
	if updated {
		err = ioutil.WriteFile(filename, data, 0777)
	}
	return
}
