package typimeta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// Handler responsible to write/read the metadata
var Handler = MetadataHandler{}

// UpdateJSON to update json metadata
func UpdateJSON(name string, v interface{}) (updated bool, err error) {
	filename := fmt.Sprintf("%s.json", name)
	var cachedData, data []byte
	cachedData, err = Handler.Read(filename)
	if os.IsNotExist(err) {
		updated = true
	} else {
		data, _ = json.Marshal(v)
		updated = (bytes.Compare(data, cachedData) != 0)
	}
	if updated {
		err = Handler.Write(filename, data)
	}
	return
}
