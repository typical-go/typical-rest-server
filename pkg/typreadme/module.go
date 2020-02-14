package typreadme

import (
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuild"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

const (
	defaultTarget   = "README.md"
	defaultTemplate = "README.tmpl"
)

// Module of readme
type Module struct {
	Target   string
	Template string
}

// New readme module
func New() *Module {
	return &Module{
		Target:   defaultTarget,
		Template: defaultTemplate,
	}
}

// WithTarget to set target
func (m *Module) WithTarget(target string) *Module {
	m.Target = target
	return m
}

// WithTemplate to set template
func (m *Module) WithTemplate(template string) *Module {
	m.Template = template
	return m
}

// BuildCommands to be shown in BuildTool
func (m *Module) BuildCommands(ctx *typbuild.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "readme",
			Usage: "Generate README Documentation",
			Action: func(c *cli.Context) (err error) {
				var (
					file *os.File
					tmpl *template.Template
					d    = ctx.Descriptor
				)
				if file, err = os.Create(m.Target); err != nil {
					return
				}
				defer file.Close()
				log.Infof("Parse template '%s'", m.Template)
				if tmpl, err = template.ParseFiles(m.Template); err != nil {
					return
				}
				obj := &ReadmeObject{
					Template:            m.Template,
					Title:               d.Name,
					Description:         d.Description,
					ApplicationCommands: appCommands(d),
					OtherBuildCommands:  otherCommands(d),
					Configs:             configs(d),
				}
				log.Infof("Apply template and write to '%s'", m.Target)
				if err = tmpl.Execute(file, obj); err != nil {
					return
				}
				return
			},
		},
	}
}

func appCommands(d *typcore.Descriptor) (infos CommandInfos) {
	if app, ok := d.App.(*typapp.App); ok {
		if app.EntryPoint() != nil {
			infos.Append(&CommandInfo{
				Snippet: d.Name,
				Usage:   "Run the application",
			})
		}
		for _, cmd := range app.AppCommands(&typapp.Context{}) {
			addCliCommandInfo(&infos, d.Name, cmd)
		}
	}
	return
}

func otherCommands(d *typcore.Descriptor) (infos CommandInfos) {
	if build, ok := d.Build.(*typbuild.Build); ok {
		for _, cmd := range build.BuildCommands(&typbuild.Context{}) {
			addCliCommandInfo(&infos, "./typicalw", cmd)
		}
	}
	return
}

func configs(d *typcore.Descriptor) (infos ConfigInfos) {
	keys, m := d.Configuration.ConfigMap()
	sort.Strings(keys)

	for _, cfg := range typcore.ConfigDetailsBy(m, keys...) {
		var required string
		if cfg.Required {
			required = "Yes"
		}
		infos.Append(&ConfigInfo{
			Name:     cfg.Name,
			Type:     cfg.Type,
			Default:  cfg.Default,
			Required: required,
		})
	}
	return
}

func addCliCommandInfo(details *CommandInfos, name string, cmd *cli.Command) {
	details.Append(&CommandInfo{
		Snippet: fmt.Sprintf("%s %s", name, cmd.Name),
		Usage:   cmd.Usage,
	})
	for _, subcmd := range cmd.Subcommands {
		details.Append(&CommandInfo{
			Snippet: fmt.Sprintf("%s %s %s", name, cmd.Name, subcmd.Name),
			Usage:   subcmd.Usage,
		})
	}
}
