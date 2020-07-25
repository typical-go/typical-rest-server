package dockerrx

import "github.com/typical-go/typical-go/pkg/typdocker"

type (
	// MySQL docker recipe
	MySQL struct {
		Version  string
		Name     string
		Image    string
		User     string
		Password string
		Port     string
	}
)

var _ typdocker.Composer = (*MySQL)(nil)

// ComposeV3 to return the recipe
func (m *MySQL) ComposeV3() (*typdocker.Recipe, error) {
	if m.Name == "" {
		m.Name = "mysql"
	}
	if m.Image == "" {
		m.Image = "mysql"
	}
	return &typdocker.Recipe{
		Version: m.Version,
		Services: typdocker.Services{
			m.Name: typdocker.Service{
				Image: m.Image,
				Environment: map[string]string{
					"MYSQL_ROOT_PASSWORD": "secret-pw", // this env is mandatory for mysql image
					"MYSQL_USER":          m.User,
					"MYSQL_PASSWORD":      m.Password,
				},
				Volumes: []string{
					m.Name + ":/var/lib/mysql",
				},
				Ports: []string{
					m.Port + ":3306",
				},
				Networks: []string{
					m.Name,
				},
				Restart: "unless-stopped",
			},
		},
		Networks: typdocker.Networks{
			m.Name: typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: typdocker.Volumes{
			m.Name: nil,
		},
	}, nil
}
