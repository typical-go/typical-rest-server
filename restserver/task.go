package restserver

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/restserver/config"
)

func taskRouteList(e *echo.Echo, cfg *config.Config) (err error) {
	m := make(map[string][]string)
	for _, route := range e.Routes() {
		if _, ok := m[route.Path]; ok {
			m[route.Path] = append(m[route.Path], route.Method)
		} else {
			m[route.Path] = []string{route.Method}
		}
	}
	fmt.Print("\n\n")
	for k, v := range m {
		fmt.Printf("/%s\t%s\n", k, v)
	}
	return
}
