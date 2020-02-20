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

// Readme module
type Readme struct {
	TargetFile   string
	TemplateFile string
	title        string
	description  string
	usages       []UsageInfo
	buildUsages  []UsageInfo
	configs      []ConfigInfo
}

// New readme module
func New() *Readme {
	return &Readme{
		TargetFile:   defaultTargetFile,
		TemplateFile: defaultTemplateFile,
	}
}

// WithTargetFile return module with new target file
func (m *Readme) WithTargetFile(targetFile string) *Readme {
	m.TargetFile = targetFile
	return m
}

// WithTemplateFile return module with new template file
func (m *Readme) WithTemplateFile(templateFile string) *Readme {
	m.TemplateFile = templateFile
	return m
}

// WithTitle return module with new title
func (m *Readme) WithTitle(title string) *Readme {
	m.title = title
	return m
}

// WithDescription return module with new description
func (m *Readme) WithDescription(description string) *Readme {
	m.description = description
	return m
}

// WithUsages return module with new usages
func (m *Readme) WithUsages(usages []UsageInfo) *Readme {
	m.usages = usages
	return m
}

// WithBuildUsages return module with new build usages
func (m *Readme) WithBuildUsages(buildUsages []UsageInfo) *Readme {
	m.buildUsages = buildUsages
	return m
}

// WithConfigs return odule with new configs
func (m *Readme) WithConfigs(configs []ConfigInfo) *Readme {
	m.configs = configs
	return m
}

// BuildCommands to be shown in BuildTool
func (m *Readme) BuildCommands(c *typbuild.Context) []*cli.Command {
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

func (m *Readme) generate(c *typbuild.Context) (err error) {
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
		Description:  m.Description(c),
		Usages:       m.Usages(c),
		BuildUsages:  m.BuildUsages(c),
		Configs:      m.Configs(c),
	})
}

// Title of readme
func (m *Readme) Title(c *typbuild.Context) string {
	if m.title == "" {
		return c.Name
	}
	return m.title
}

// Description of readme
func (m *Readme) Description(c *typbuild.Context) string {
	if m.description == "" {
		return c.Description
	}
	return m.description
}

// Usages of readme
func (m *Readme) Usages(c *typbuild.Context) (infos []UsageInfo) {
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
func (m *Readme) BuildUsages(c *typbuild.Context) (infos []UsageInfo) {
	if len(m.buildUsages) < 1 {
		if build, ok := c.BuildTool.(*typbuildtool.BuildTool); ok {
			for _, cmd := range build.BuildCommands(&typbuild.Context{}) {
				infos = append(infos, usageInfos("./typicalw", cmd)...)
			}
		}
		return
	}
	return m.buildUsages
}

// Configs of readme
func (m *Readme) Configs(c *typbuild.Context) (infos []ConfigInfo) {
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
