package typigen

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

type SampleConfigType1 struct{}
type SampleConfigType2 struct{}

func TestConstructConfig(t *testing.T) {
	mainConfig, configConstructors := constructConfig(typictx.Context{
		Modules: []*typictx.Module{
			{ConfigPrefix: "APP", ConfigSpec: &SampleConfigType1{}, Name: "Application configuration"},
			{ConfigPrefix: "PG", ConfigSpec: &SampleConfigType2{}, Name: "Postgres configuration"},
		},
	})

	require.Equal(t, "type Config struct{\nApp *typigen.SampleConfigType1\nPg *typigen.SampleConfigType2\n}\n", mainConfig.String())
	require.Equal(t, 3, len(configConstructors))
	require.Equal(t, "func () (*Config,error){ \nvar cfg Config\nerr := envconfig.Process(\"\", &cfg)\nreturn &cfg, err\n}", configConstructors[0].String())
	require.Equal(t, "func (cfg *Config,) (*typigen.SampleConfigType1){ \nreturn cfg.App\n}", configConstructors[1].String())
	require.Equal(t, "func (cfg *Config,) (*typigen.SampleConfigType2){ \nreturn cfg.Pg\n}", configConstructors[2].String())
}
