package typictx

import "github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"

// Configurations return config list
func Configurations(ctx *Context) (cfgs []typiobj.Configuration) {
	if configurer, ok := ctx.Application.(typiobj.Configurer); ok {
		cfgs = append(cfgs, configurer.Configure())
	}
	cfgs = append(cfgs, ctx.Modules.Configurations()...)
	return
}
