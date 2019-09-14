package prebuilder

import (
	"encoding/json"
	"io/ioutil"
)

func writeCache(filename string, val interface{}) error {
	data, _ := json.MarshalIndent(val, "", "    ")
	return ioutil.WriteFile(".typical/"+filename, data, 0644)
}
