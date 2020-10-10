package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app"
	_ "github.com/typical-go/typical-rest-server/internal/generated/typical"
)

func main() {
	fmt.Printf("%s %s\n", typgo.ProjectName, typgo.ProjectVersion)
	if err := typapp.Run(app.Start); err != nil {
		logrus.Fatal(err.Error())
	}
}
