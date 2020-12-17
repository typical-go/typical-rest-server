package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	a "github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary/service"
	b "github.com/typical-go/typical-rest-server/internal/app/domain/mymusic/service"
	c "github.com/typical-go/typical-rest-server/internal/app/infra"
)

func init() {
	typapp.Provide("", a.NewBookSvc)
	typapp.Provide("", b.NewSongSvc)
	typapp.Provide("", c.NewCacheStore)
	typapp.Provide("", c.NewDatabases)
	typapp.Provide("", c.NewServer)
}
