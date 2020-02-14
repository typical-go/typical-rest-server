package typreadme

import (
	"fmt"
	"html/template"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

const (
	defaultTargetFile   = "README.md"
	defaultTemplateFile = "README.tmpl"
)

// Module of readme
type Module struct {
	TargetFile   string
	TemplateFile string
}

// New readme module
func New() *Module {
	return &Module{
		TargetFile:   defaultTargetFile,
		TemplateFile: defaultTemplateFile,
	}
}

// WithTargetFile to set target file
func (m *Module) WithTargetFile(targetFile string) *Module {
	m.TargetFile = targetFile
	return m
}

// WithTemplateFile to set template
func (m *Module) WithTemplateFile(templateFile string) *Module {
	m.TemplateFile = templateFile
	return m
}

// BuildCommands to be shown in BuildTool
func (m *Module) BuildCommands(c *typbuild.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "readme",
			Usage: "Generate README Documentation",
			Action: func(cliCtx *cli.Context) (err error) {
				return m.generate(c)
			},
		},
	}
}

func (m *Module) generate(c *typbuild.Context) (err error) {
	var (
		file *os.File
		tmpl *template.Template
	)
	if file, err = os.Create(m.TargetFile); err != nil {
		return
	}
	defer file.Close()
	log.Infof("Parse template '%s'", m.TemplateFile)
	if tmpl, err = template.ParseFiles(m.TemplateFile); err != nil {
		return
	}
	log.Infof("Apply template and write to '%s'", m.TargetFile)
	return tmpl.Execute(file, &Object{
		TemplateFile:  m.TemplateFile,
		Title:         c.Name,
		Description:   c.Description,
		Usages:        m.appCommands(c),
		BuildCommands: m.otherCommands(c),
		Configs:       m.configs(c),
	})
}

func (m *Module) appCommands(c *typbuild.Context) (infos []CommandInfo) {
	if app, ok := c.App.(*typapp.App); ok {
		if app.EntryPoint() != nil {
			infos = append(infos, CommandInfo{
				Command:     c.Name,
				Description: "Run the application",
			})
		}
		for _, cmd := range app.AppCommands(&typapp.Context{}) {
			infos = append(infos, commandInfos(c.Name, cmd)...)
		}
	}
	return
}

func (m *Module) otherCommands(c *typbuild.Context) (infos []CommandInfo) {
	if build, ok := c.Build.(*typbuild.Build); ok {
		for _, cmd := range build.BuildCommands(&typbuild.Context{}) {
			infos = append(infos, commandInfos("./typicalw", cmd)...)
		}
	}
	return
}

func (m *Module) configs(c *typbuild.Context) (infos []ConfigInfo) {
	keys, cfgmap := c.Configuration.ConfigMap()
	sort.Strings(keys)

	for _, cfg := range typcore.ConfigDetailsBy(cfgmap, keys...) {
		var required string
		if cfg.Required {
			required = "Yes"
		}
		infos = append(infos, ConfigInfo{
			Name:     cfg.Name,
			Type:     cfg.Type,
			Default:  cfg.Default,
			Required: required,
		})
	}
	return
}

func commandInfos(name string, cmd *cli.Command) (details []CommandInfo) {
	details = append(details, CommandInfo{
		Command:     fmt.Sprintf("%s %s", name, cmd.Name),
		Description: cmd.Usage,
	})
	for _, subcmd := range cmd.Subcommands {
		details = append(details, CommandInfo{
			Command:     fmt.Sprintf("%s %s %s", name, cmd.Name, subcmd.Name),
			Description: subcmd.Usage,
		})
	}
	return
}
