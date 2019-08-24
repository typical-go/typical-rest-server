package typitask

import (
	"os/exec"
	"strings"
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
	configTemplate = `| Key | Type | Default | Required | Description |	
|---|---|---|---|---|{{range .}}
|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|{{end}}`
)

// GenerateReadme for generate typical applical readme
func GenerateReadme(a *typictx.ActionContext) (err error) {

	configDoc := ConfigDoc{Context: a.Context}

	readme0 := readme.DefaultReadme().
		SetTitle(a.Name).
		SetDescription(a.Description).
		SetSection("Configuration", configDoc.Section)

	log.Infof("Generate README.md")
	err = readme0.OutputToFile("README.md")
	if err != nil {
		return
	}

	if !a.Cli.Bool("no-commit") {
		gitRepo, _ := git.PlainOpen(".")
		worktree, _ := gitRepo.Worktree()
		status, _ := worktree.Status()

		fileStatus := status.File("README.md")
		if fileStatus.Worktree != git.Modified {
			log.Info("No change at README.md")
			return
		}

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
		log.Info("Push change at README.md")
		return exec.Command("git", "push").Run()
	}
	return
}

type ConfigDoc struct {
	*typictx.Context
}

func (d ConfigDoc) Section(md *typirecipe.Markdown) (err error) {
	for i, acc := range d.Context.ConfigAccessors() {
		name := acc.GetName()
		if name != "" {
			if i > 0 {

			}
			md.Heading3(name)
		}

		var builder strings.Builder
		envconfig.Usagef(acc.GetConfigPrefix(), acc.GetConfigSpec(), &builder, configTemplate)
		md.Writeln(builder.String())
	}
	return
}
