package typcfg_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/typical-go/typical-go/pkg/oskit"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
)

func TestCfgAnnotation_Annotate(t *testing.T) {
	typgo.ProjectPkg = "github.com/user/project"
	defer os.RemoveAll("internal")

	defer typgo.PatchBash([]*typgo.RunExpectation{})(t)

	var out strings.Builder
	defer oskit.PatchStdout(&out)()

	EnvconfigAnnotation := &typcfg.EnvconfigAnnotation{}
	c := &typast.Context{
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typast.Summary{
			Annots: []*typast.Annot{
				{
					TagName: "@envconfig",
					Decl: &typast.Decl{
						File: typast.File{Package: "mypkg", Path: "pkg/file.go"},
						Type: &typast.StructDecl{
							TypeDecl: typast.TypeDecl{Name: "SomeSample"},
							Fields: []*typast.Field{
								{Names: []string{"SomeField1"}, Type: "string", StructTag: `default:"some-text"`},
								{Names: []string{"SomeField2"}, Type: "int", StructTag: `default:"9876"`},
							},
						},
					},
				},
			},
		},
	}

	require.NoError(t, EnvconfigAnnotation.Annotate(c))

	b, _ := ioutil.ReadFile("internal/generated/envcfg/envcfg.go")
	require.Equal(t, `package envcfg

/* DO NOT EDIT. This file generated due to '@envconfig' annotation */

import (
	 "fmt"
	 "github.com/kelseyhightower/envconfig"
	 "github.com/typical-go/typical-go/pkg/typapp"
	a "github.com/user/project/pkg"
)

func init() { 
	typapp.Provide("",LoadSomeSample)
}

// LoadSomeSample load env to new instance of SomeSample
func LoadSomeSample() (*a.SomeSample, error) {
	var cfg a.SomeSample
	prefix := "SOMESAMPLE"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}
`, string(b))

	require.Equal(t, "Generate @envconfig to internal/generated/envcfg/envcfg.go\n", out.String())

}

func TestCfgAnnotation_Annotate_GenerateDotEnvAndUsageDoc(t *testing.T) {
	var out strings.Builder

	defer typgo.PatchBash(nil)(t)
	defer oskit.PatchStdout(&out)()
	defer os.Clearenv()
	defer os.RemoveAll("folder")

	a := &typcfg.EnvconfigAnnotation{
		Target:   "folder/some-target",
		Template: "some-template",
		DotEnv:   ".env33",
		UsageDoc: "some-usage.md",
	}
	c := &typast.Context{
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typast.Summary{Annots: []*typast.Annot{
			{
				TagName:  "@envconfig",
				TagParam: `ctor:"ctor1" prefix:"SS"`,
				Decl: &typast.Decl{
					File: typast.File{Package: "mypkg"},
					Type: &typast.StructDecl{
						TypeDecl: typast.TypeDecl{Name: "SomeSample"},
						Fields: []*typast.Field{
							{Names: []string{"SomeField1"}, Type: "string", StructTag: `default:"some-text"`},
							{Names: []string{"SomeField2"}, Type: "int", StructTag: `default:"9876"`},
						},
					},
				},
			},
		}},
	}

	require.NoError(t, a.Annotate(c))
	defer os.Remove(a.Target)
	defer os.Remove(a.DotEnv)
	defer os.Remove(a.UsageDoc)

	b, _ := ioutil.ReadFile(a.Target)
	require.Equal(t, `some-template`, string(b))

	b, _ = ioutil.ReadFile(a.DotEnv)
	require.Equal(t, "SS_SOMEFIELD1=some-text\nSS_SOMEFIELD2=9876\n", string(b))
	require.Equal(t, "some-text", os.Getenv("SS_SOMEFIELD1"))
	require.Equal(t, "9876", os.Getenv("SS_SOMEFIELD2"))

	require.Equal(t, "Generate @envconfig to folder/some-target\nNew keys added in '.env33': SS_SOMEFIELD1 SS_SOMEFIELD2\nGenerate 'some-usage.md'\n", out.String())
}

func TestCfgAnnotation_Annotate_Predefined(t *testing.T) {
	defer typgo.PatchBash(nil)(t)
	defer os.RemoveAll("predefined")

	EnvconfigAnnotation := &typcfg.EnvconfigAnnotation{
		TagName:  "@some-tag",
		Template: "some-template",
		Target:   "predefined/cfg-target",
	}
	c := &typast.Context{
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typast.Summary{
			Annots: []*typast.Annot{
				{
					TagName: "@some-tag",
					Decl: &typast.Decl{
						File: typast.File{Package: "mypkg"},
						Type: &typast.StructDecl{
							TypeDecl: typast.TypeDecl{Name: "SomeSample"},
							Fields:   []*typast.Field{},
						},
					},
				},
			},
		},
	}
	require.NoError(t, EnvconfigAnnotation.Annotate(c))

	b, _ := ioutil.ReadFile("predefined/cfg-target")
	require.Equal(t, `some-template`, string(b))
}

func TestCfgAnnotation_Annotate_RemoveTargetWhenNoAnnotation(t *testing.T) {
	target := "target1"
	defer os.Remove(target)
	ioutil.WriteFile(target, []byte("some-content"), 0777)
	c := &typast.Context{
		Context: &typgo.Context{},
		Summary: &typast.Summary{},
	}

	EnvconfigAnnotation := &typcfg.EnvconfigAnnotation{Target: target}
	require.NoError(t, EnvconfigAnnotation.Annotate(c))
	_, err := os.Stat(target)
	require.True(t, os.IsNotExist(err))
}

func TestCreateField(t *testing.T) {
	testnames := []struct {
		TestName string
		Prefix   string
		Field    *typast.Field
		Expected *typcfg.Field
	}{
		{
			Prefix:   "APP",
			Field:    &typast.Field{Names: []string{"Address"}},
			Expected: &typcfg.Field{Key: "APP_ADDRESS"},
		},
		{
			Prefix: "APP",
			Field: &typast.Field{
				Names:     []string{"some-name"},
				StructTag: reflect.StructTag(`envconfig:"ADDRESS" default:"some-address" required:"true"`),
			},
			Expected: &typcfg.Field{Key: "APP_ADDRESS", Default: "some-address", Required: true},
		},
	}
	for _, tt := range testnames {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.Expected, typcfg.CreateField(tt.Prefix, tt.Field))
		})
	}
}
