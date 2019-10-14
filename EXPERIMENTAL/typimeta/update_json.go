package typimeta

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

// UpdateJSON to update json metadata
func UpdateJSON(filename string, v interface{}) (updated bool, err error) {
	var cachedData, data []byte
	cachedData, err = ioutil.ReadFile(filename)
	data, _ = json.Marshal(v)
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
