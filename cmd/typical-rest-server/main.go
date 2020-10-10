package main

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/internal/app/infra/log"
	_ "github.com/typical-go/typical-rest-server/internal/generated/typical"
)

func main() {
	fmt.Printf("%s %s\n", typgo.ProjectName, typgo.ProjectVersion)
	if err := typapp.Run(app.Start); err != nil {
		log.Fatal(err.Error())
	}
}
