package dockerrx

import (
	"errors"
	"os"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// MySQLWithEnv same with MySQL but with env parameter
type MySQLWithEnv struct {
	Version     string
	Name        string
	Image       string
	UserEnv     string
	PasswordEnv string
	PortEnv     string
}

var _ (typdocker.Composer) = (*MySQLWithEnv)(nil)

// Compose for docker-compose
func (m *MySQLWithEnv) Compose() (*typdocker.Recipe, error) {
	if m.PasswordEnv == "" {
		return nil, errors.New("mysql-with-env: missing PasswordEnv")
	}
	if m.PortEnv == "" {
		return nil, errors.New("mysql-with-env: missing PortEnv")
	}
	if m.UserEnv == "" {
		return nil, errors.New("mysql-with-env: missing UserEnv")
	}
	mysql := &MySQL{
		Version:  m.Version,
		Name:     m.Name,
		Image:    m.Image,
		User:     os.Getenv(m.UserEnv),
		Password: os.Getenv(m.PasswordEnv),
		Port:     os.Getenv(m.PortEnv),
	}
	return mysql.Compose()
}
