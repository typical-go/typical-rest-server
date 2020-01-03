package typreadme

import (
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (r *readme) generateCmd() *cli.Command {
	return &cli.Command{
		Name:  "readme",
		Usage: "Generate README.md",
		Action: func(c *cli.Context) (err error) {
			var file *os.File
			if file, err = os.Create("README.md"); err != nil {
				return
			}
			defer file.Close()
			var tmpl *template.Template
			if tmpl, err = template.ParseFiles("README.tmpl"); err != nil {
				return
			}
			if err = tmpl.Execute(file, r.readmeObj()); err != nil {
				return
			}
			return
		},
	}
}

func (r *readme) readmeObj() *ReadmeObject {
	return &ReadmeObject{
		Title:               r.Name,
		Description:         r.Description,
		ApplicationCommands: r.appCommands(),
		OtherBuildCommands:  r.otherCommands(),
		Configs:             r.configs(),
	}
}

func (r *readme) appCommands() (infos CommandInfos) {
	appName := strcase.ToKebab(r.Name) // TODO: use typenv instead
	if typcore.IsActionable(r.AppModule) {
		infos.Append(&CommandInfo{
			Snippet: appName,
			Usage:   "Run the application",
		})
	}
	if commander, ok := r.AppModule.(typcore.AppCommander); ok {
		for _, cmd := range commander.AppCommands(&typcore.Context{}) {
			addCliCommandInfo(&infos, appName, cmd)
		}
	}
	return
}

func (r *readme) otherCommands() (infos CommandInfos) {
	for _, cmd := range typbuildtool.BuildCommands(r.ProjectDescriptor) {
		addCliCommandInfo(&infos, "./typicalw", cmd)
	}
	return
}

func (r *readme) configs() (infos ConfigInfos) {
	keys, configMap := typcore.CreateConfigMap(r.ProjectDescriptor)
	sort.Strings(keys)
	for _, cfg := range configMap.ValueBy(keys...) {
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
