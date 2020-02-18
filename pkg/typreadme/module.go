package typreadme

import (
	"fmt"
	"html/template"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
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
	title        string
	description  string
	usages       []UsageInfo
	buildUsages  []UsageInfo
	configs      []ConfigInfo
}

// New readme module
func New() *Module {
	return &Module{
		TargetFile:   defaultTargetFile,
		TemplateFile: defaultTemplateFile,
	}
}

// WithTargetFile return module with new target file
func (m *Module) WithTargetFile(targetFile string) *Module {
	m.TargetFile = targetFile
	return m
}

// WithTemplateFile return module with new template file
func (m *Module) WithTemplateFile(templateFile string) *Module {
	m.TemplateFile = templateFile
	return m
}

// WithTitle return module with new title
func (m *Module) WithTitle(title string) *Module {
	m.title = title
	return m
}

// WithDescription return module with new description
func (m *Module) WithDescription(description string) *Module {
	m.description = description
	return m
}

// WithUsages return module with new usages
func (m *Module) WithUsages(usages []UsageInfo) *Module {
	m.usages = usages
	return m
}

// WithBuildUsages return module with new build usages
func (m *Module) WithBuildUsages(buildUsages []UsageInfo) *Module {
	m.buildUsages = buildUsages
	return m
}

// WithConfigs return odule with new configs
func (m *Module) WithConfigs(configs []ConfigInfo) *Module {
	m.configs = configs
	return m
}

// buildUsages  []UsageInfo
// 	configs      []ConfigInfo

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
		TemplateFile: m.TemplateFile,
		Title:        m.Title(c),
		Description:  m.description,
		Usages:       m.Usages(c),
		BuildUsages:  m.BuildUsages(c),
		Configs:      m.Configs(c),
	})
}

// Title of readme
func (m *Module) Title(c *typbuild.Context) string {
	if m.title == "" {
		return c.Name
	}
	return m.title
}

// Usages of readme
func (m *Module) Usages(c *typbuild.Context) (infos []UsageInfo) {
	if len(m.usages) < 1 {
		if app, ok := c.App.(*typapp.App); ok {
			if app.EntryPoint() != nil {
				infos = append(infos, UsageInfo{
					Usage:       c.Name,
					Description: "Run the application",
				})
			}
			for _, cmd := range app.AppCommands(&typapp.Context{}) {
				infos = append(infos, usageInfos(c.Name, cmd)...)
			}
		}
		return
	}
	return m.usages
}

// BuildUsages of readme
func (m *Module) BuildUsages(c *typbuild.Context) (infos []UsageInfo) {
	if len(m.buildUsages) < 1 {
		if build, ok := c.BuildTool.(*typbuildtool.Build); ok {
			for _, cmd := range build.BuildCommands(&typbuild.Context{}) {
				infos = append(infos, usageInfos("./typicalw", cmd)...)
			}
		}
		return
	}
	return m.buildUsages
}

// Configs of readme
func (m *Module) Configs(c *typbuild.Context) (infos []ConfigInfo) {
	if len(m.configs) < 1 {
		store := c.Configuration.Store()
		keys := store.Keys()
		sort.Strings(keys)
		for _, cfg := range store.Fields(keys...) {
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
	return m.configs
}

func usageInfos(name string, cmd *cli.Command) (details []UsageInfo) {
	details = append(details, UsageInfo{
		Usage:       fmt.Sprintf("%s %s", name, cmd.Name),
		Description: cmd.Usage,
	})
	for _, subcmd := range cmd.Subcommands {
		details = append(details, UsageInfo{
			Usage:       fmt.Sprintf("%s %s %s", name, cmd.Name, subcmd.Name),
			Description: subcmd.Usage,
		})
	}
	return
}
