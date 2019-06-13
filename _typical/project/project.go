package project

import (
	"bytes"
	"log"

	"github.com/BurntSushi/toml"
)

func ContextDetail() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(Context); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
