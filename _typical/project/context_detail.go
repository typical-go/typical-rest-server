package project

import (
	"bytes"
	"log"

	"github.com/BurntSushi/toml"
)

// ContextDetail return context detail string
func ContextDetail() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(Ctx); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
