package projctx

import (
	"bytes"
	"log"

	"github.com/BurntSushi/toml"
)

var ctx Context

func init() {
	_, err := toml.DecodeFile(".typical/_context.toml", &ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func Name() string {
	return ctx.Name
}

func Description() string {
	return ctx.Description
}

func Version() string {
	return ctx.Version
}

func String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(ctx); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
