package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	a "github.com/typical-go/typical-rest-server/internal/app/infra"
	b "github.com/typical-go/typical-rest-server/internal/app/repo"
	c "github.com/typical-go/typical-rest-server/internal/app/service"
)

func init() {
	typapp.Provide("", a.NewCacheStore)
	typapp.Provide("", a.NewDatabases)
	typapp.Provide("", a.NewEcho)
	typapp.Provide("", b.NewBookRepo)
	typapp.Provide("", c.NewBookSvc)
}
