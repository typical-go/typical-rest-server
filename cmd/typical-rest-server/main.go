package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app"

	_ "github.com/typical-go/typical-rest-server/internal/generated/ctor"
	_ "github.com/typical-go/typical-rest-server/internal/generated/envcfg"
)

func main() {
	fmt.Printf("%s %s\n", typgo.ProjectName, typgo.ProjectVersion)
	if err := typapp.StartApp(app.Start, app.Shutdown); err != nil {
		logrus.Fatal(err.Error())
	}

}
