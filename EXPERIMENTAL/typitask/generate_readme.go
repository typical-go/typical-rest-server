package typitask

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-gitconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/readme"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	configTemplate = `
| Key | Type | Default | Required | Description |	
|---|---|---|---|---|	
{{range .}}|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|	
{{end}}`
)

// GenerateReadme for generate typical applical readme
func GenerateReadme(ctx *typictx.ActionContext) (err error) {
	recipe := readme.Recipe{
		Title:       ctx.Typical.Name,
		Description: ctx.Typical.Description,
		Sections: []readme.Section{
			{Title: "Getting Started", Content: gettingStartedInstruction()},
			{Title: "Usage", Content: usageInstruction()},
			{Title: "Build Tool", Content: buildToolInstruction()},
			{Title: "Make a release", Content: releaseInstruction()},
			{Title: "Configurations", Content: configDoc(ctx.Typical)},
		},
	}

	log.Infof("Generate README.md")
	err = recipe.WriteToFile("README.md")
	if err != nil {
		return
	}

	token := ctx.Cli.String("github-token")
	if token != "" {
		gitRepo, _ := git.PlainOpen(".")
		worktree, _ := gitRepo.Worktree()
		status, _ := worktree.Status()

		fileStatus := status.File("README.md")
		if fileStatus.Worktree == git.Modified {
			user, _ := gitconfig.Username()
			email, _ := gitconfig.Email()

			worktree.Add("README.md")
			_, err = worktree.Commit("Generate latest README.md", &git.CommitOptions{
				Author: &object.Signature{
					Name:  user,
					Email: email,
					When:  time.Now(),
				},
			})
			if err != nil {
				return
			}

			return exec.Command("git", "push").Run()
		}
	}

	return
}

func configDoc(ctx typictx.Context) string {
	buf := new(bytes.Buffer)

	for i, cfg := range ctx.Configurations {
		if i > 0 {
			buf.WriteString("\n")
		}

		buf.WriteString(cfg.Description)
		buf.WriteString("\n")
		envconfig.Usagef(cfg.Prefix, cfg.Spec, buf, configTemplate)
	}

	return buf.String()
}

func gettingStartedInstruction() string {
	var md typirecipe.Markdown
	md.Writeln("This is intruction to start working with the project:")
	md.OrderedList(
		"Install [Go](https://golang.org/doc/install) or using homebrew if you're using macOS `brew install go`",
	)

	return md.String()
}

func usageInstruction() string {
	return `There is no specific requirement to run the application. `
}

func buildToolInstruction() string {
	return "Use `./typicalw` to execute development task"
}

func releaseInstruction() string {
	return "Use `./typicalw release -github-token=[TOKEN]` to make the release. You can found the release in `release` folder or github release page"
}
