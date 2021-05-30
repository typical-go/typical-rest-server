package envcfg

/* DO NOT EDIT. This file generated due to '@envconfig' annotation */

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-go/pkg/typapp"
	a "github.com/typical-go/typical-rest-server/internal/app/infra"
)

func init() {
	typapp.Provide("", LoadAppCfg)
	typapp.Provide("", LoadCacheCfg)
	typapp.Provide("pg", LoadPgDatabaseCfg)
}

// LoadAppCfg load env to new instance of AppCfg
func LoadAppCfg() (*a.AppCfg, error) {
	var cfg a.AppCfg
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadCacheCfg load env to new instance of CacheCfg
func LoadCacheCfg() (*a.CacheCfg, error) {
	var cfg a.CacheCfg
	prefix := "CACHE"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadPgDatabaseCfg load env to new instance of DatabaseCfg
func LoadPgDatabaseCfg() (*a.DatabaseCfg, error) {
	var cfg a.DatabaseCfg
	prefix := "PG"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}
