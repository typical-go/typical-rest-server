package dockerrx

import (
	"errors"
	"os"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

type (
	// MySQLWithEnv same with MySQL but with env parameter
	MySQLWithEnv struct {
		Version     string
		Name        string
		Image       string
		UserEnv     string
		PasswordEnv string
		PortEnv     string
	}
)

var _ (typdocker.Composer) = (*MySQLWithEnv)(nil)

// ComposeV3 for docker-compose
func (m *MySQLWithEnv) ComposeV3() (*typdocker.Recipe, error) {
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
	return mysql.ComposeV3()
}
