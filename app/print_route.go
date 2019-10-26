package app

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/config"
)

// DryRun just run
func DryRun(e *echo.Echo, cfg *config.Config) (err error) {
	data, _ := json.MarshalIndent(e.Routes(), "", "  ")
	fmt.Println(string(data))
	return
}
