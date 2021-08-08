package typcfg_test

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgen"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

func TestCfgAnnotation_Annotate(t *testing.T) {
	typgo.ProjectPkg = "github.com/user/project"
	defer os.RemoveAll("internal")

	EnvconfigAnnot := &typcfg.EnvconfigAnnot{}
	var out strings.Builder
	c := &typgo.Context{Logger: typgo.Logger{Stdout: &out}}
	defer c.PatchBash([]*typgo.MockBash{})(t)

	directives := []*typgen.Directive{
		{
			TagName: "@envconfig",
			Decl: &typgen.Decl{
				File: typgen.File{Package: "mypkg", Path: "pkg/file.go"},
				Type: &typgen.StructDecl{
					TypeDecl: typgen.TypeDecl{Name: "SomeSample"},
					Fields: []*typgen.Field{
						{Names: []string{"SomeField1"}, Type: "string", StructTag: `default:"some-text"`},
						{Names: []string{"SomeField2"}, Type: "int", StructTag: `default:"9876"`},
					},
				},
			},
		},
	}

	require.NoError(t, EnvconfigAnnot.Process(c, directives))

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

	require.Equal(t, "> Generate @envconfig to internal/generated/envcfg/envcfg.go\n> go build -o /bin/goimports golang.org/x/tools/cmd/goimports\n", out.String())

}

func TestCfgAnnotation_Annotate_GenerateDotEnvAndUsageDoc(t *testing.T) {

	defer os.Clearenv()
	defer os.RemoveAll("folder")

	a := &typcfg.EnvconfigAnnot{
		Target:    "folder/some-target",
		Template:  "some-template",
		GenDotEnv: ".env33",
		GenDoc:    "some-usage.md",
	}

	var out strings.Builder
	c := &typgo.Context{
		Context:    cli.NewContext(nil, &flag.FlagSet{}, nil),
		Descriptor: &typgo.Descriptor{},
		Logger:     typgo.Logger{Stdout: &out},
	}
	defer c.PatchBash(nil)(t)

	directives := []*typgen.Directive{
		{
			TagName:  "@envconfig",
			TagParam: `ctor:"ctor1" prefix:"SS"`,
			Decl: &typgen.Decl{
				File: typgen.File{Package: "mypkg"},
				Type: &typgen.StructDecl{
					TypeDecl: typgen.TypeDecl{Name: "SomeSample"},
					Fields: []*typgen.Field{
						{Names: []string{"SomeField1"}, Type: "string", StructTag: `default:"some-text"`},
						{Names: []string{"SomeField2"}, Type: "int", StructTag: `default:"9876"`},
					},
				},
			},
		},
		{
			TagName:  "@envconfig",
			TagParam: `prefix:"-"`,
			Decl: &typgen.Decl{
				File: typgen.File{Package: "mypkg"},
				Type: &typgen.StructDecl{
					TypeDecl: typgen.TypeDecl{Name: "SomeSample"},
					Fields: []*typgen.Field{
						{Names: []string{"SomeField3"}, Type: "string", StructTag: `default:"some-text"`},
						{Names: []string{"SomeField4"}, Type: "int", StructTag: `default:"9876"`},
					},
				},
			},
		},
		{
			TagName:  "@envconfig",
			TagParam: `prefix:"_"`,
			Decl: &typgen.Decl{
				File: typgen.File{Package: "mypkg"},
				Type: &typgen.StructDecl{
					TypeDecl: typgen.TypeDecl{Name: "SomeSample"},
					Fields: []*typgen.Field{
						{Names: []string{"SomeField5"}, Type: "string", StructTag: `default:"some-text"`},
						{Names: []string{"SomeField6"}, Type: "int", StructTag: `default:"9876"`},
					},
				},
			},
		},
	}

	require.NoError(t, a.Process(c, directives))
	defer os.Remove(a.Target)
	defer os.Remove(a.GenDotEnv)
	defer os.Remove(a.GenDoc)

	b, _ := ioutil.ReadFile(a.Target)
	require.Equal(t, `some-template`, string(b))

	b, _ = ioutil.ReadFile(a.GenDotEnv)
	require.Equal(t, `SOMEFIELD3=some-text
SOMEFIELD4=9876
SOMEFIELD5=some-text
SOMEFIELD6=9876
SS_SOMEFIELD1=some-text
SS_SOMEFIELD2=9876
`, string(b))
	require.Equal(t, "some-text", os.Getenv("SS_SOMEFIELD1"))
	require.Equal(t, "9876", os.Getenv("SS_SOMEFIELD2"))
}

func TestCfgAnnotation_Annotate_Predefined(t *testing.T) {

	defer os.RemoveAll("predefined")

	EnvconfigAnnot := &typcfg.EnvconfigAnnot{
		TagName:  "@some-tag",
		Template: "some-template",
		Target:   "predefined/cfg-target",
	}
	c := &typgo.Context{}
	defer c.PatchBash(nil)(t)

	directives := []*typgen.Directive{
		{
			TagName: "@some-tag",
			Decl: &typgen.Decl{
				File: typgen.File{Package: "mypkg"},
				Type: &typgen.StructDecl{
					TypeDecl: typgen.TypeDecl{Name: "SomeSample"},
					Fields:   []*typgen.Field{},
				},
			},
		},
	}

	require.NoError(t, EnvconfigAnnot.Process(c, directives))

	b, _ := ioutil.ReadFile("predefined/cfg-target")
	require.Equal(t, `some-template`, string(b))
}

func TestCfgAnnotation_Annotate_RemoveTargetWhenNoAnnotation(t *testing.T) {
	target := "target1"
	defer os.Remove(target)
	ioutil.WriteFile(target, []byte("some-content"), 0777)

	EnvconfigAnnot := &typcfg.EnvconfigAnnot{Target: target}
	require.NoError(t, EnvconfigAnnot.Process(nil, nil))
	_, err := os.Stat(target)
	require.True(t, os.IsNotExist(err))
}

func TestCreateField(t *testing.T) {
	testnames := []struct {
		TestName string
		Prefix   string
		Field    *typgen.Field
		Expected *typcfg.Field
	}{
		{
			Prefix:   "APP",
			Field:    &typgen.Field{Names: []string{"Address"}},
			Expected: &typcfg.Field{Key: "APP_ADDRESS"},
		},
		{
			Prefix: "APP",
			Field: &typgen.Field{
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
