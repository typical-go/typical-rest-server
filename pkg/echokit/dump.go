package echokit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

// DumpEcho dump route in echo server
func DumpEcho(e *echo.Echo) []string {
	m := make(map[string][]string)
	for _, r := range e.Routes() {
		path := r.Path
		method := r.Method
		if _, ok := m[r.Path]; ok {
			m[path] = append(m[path], method)
		} else {
			m[path] = []string{method}
		}
	}

	var paths []string
	for k, v := range m {
		sort.Strings(v)
		paths = append(paths, fmt.Sprintf("%s\t%s", k, strings.Join(v, ",")))
	}
	sort.Strings(paths)
	return paths
}
