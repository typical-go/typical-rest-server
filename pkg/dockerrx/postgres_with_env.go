package dockerrx

import (
	"errors"
	"os"

	"github.com/typical-go/typical-go/pkg/typdocker"
)

// PostgresWithEnv same with postgres with env parameter
type PostgresWithEnv struct {
	Version     string
	Name        string
	Image       string
	UserEnv     string
	PasswordEnv string
	PortEnv     string
}

var _ (typdocker.Composer) = (*PostgresWithEnv)(nil)

// Compose for docker-compose
func (p *PostgresWithEnv) Compose() (*typdocker.Recipe, error) {
	if p.PasswordEnv == "" {
		return nil, errors.New("pg-with-env: missing PasswordEnv")
	}
	if p.PortEnv == "" {
		return nil, errors.New("pg-with-env: missing PortEnv")
	}
	if p.UserEnv == "" {
		return nil, errors.New("pg-with-env: missing UserEnv")
	}
	pg := &Postgres{
		Version:  p.Version,
		Name:     p.Name,
		Image:    p.Image,
		User:     os.Getenv(p.UserEnv),
		Password: os.Getenv(p.PasswordEnv),
		Port:     os.Getenv(p.PortEnv),
	}
	return pg.Compose()
}
