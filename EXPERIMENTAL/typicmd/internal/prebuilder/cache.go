package prebuilder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func writeCache(name string, val interface{}) error {
	data, _ := json.MarshalIndent(val, "", "    ")
	return ioutil.WriteFile(cacheFile(name), data, 0644)
}

func readCache(name string, v interface{}) (err error) {
	file, err := os.Open(cacheFile(name))
	if err != nil {
		return
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	return json.Unmarshal(data, v)
}

func cacheFile(name string) string {
	return fmt.Sprintf(".typical/%s.json", name)
}
