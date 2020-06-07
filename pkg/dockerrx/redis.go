package dockerrx

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// Redis docker recipe
type Redis struct {
	Version  string
	Name     string
	Image    string
	Password string
	Port     string
}

var _ typdocker.Composer = (*Redis)(nil)

// Compose to return redis recipe
func (r *Redis) Compose() (*typdocker.Recipe, error) {
	return &typdocker.Recipe{
		Version: r.Version,
		Services: typdocker.Services{
			r.Name: typdocker.Service{
				Image:    r.Image,
				Command:  fmt.Sprintf(`redis-server --requirepass "%s"`, r.Password),
				Ports:    []string{fmt.Sprintf("%s:6379", r.Port)},
				Networks: []string{r.Name},
				Volumes:  []string{fmt.Sprintf("%s:/data", r.Name)},
			},
		},
		Networks: typdocker.Networks{
			r.Name: nil,
		},
		Volumes: typdocker.Volumes{
			r.Name: nil,
		},
	}, nil
}
