package buildtool

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/buildtool/releaser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
	readmeFile = "README.md"
)

func commands(c *typictx.Context) (cmds []*typictx.Command) {
	cmds = []*typictx.Command{
		{
			Name:       "build",
			ShortName:  "b",
			Usage:      "Build the binary",
			ActionFunc: buildBinary,
		},
		{
			Name:       "clean",
			ShortName:  "c",
			Usage:      "Clean the project from generated file during build time",
			ActionFunc: cleanProject,
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the binary",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-build",
					Usage: "Run the binary without build",
				},
			},
			ActionFunc: runBinary,
		},
		{
			Name:       "test",
			ShortName:  "t",
			Usage:      "Run the testing",
			ActionFunc: runTesting,
		},
		{
			Name:  "release",
			Usage: "Release the distribution",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-test",
					Usage: "Release without run automated test",
				},
				cli.BoolFlag{
					Name:  "no-github",
					Usage: "Release without create github release",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Release by passed all validation",
				},
				cli.BoolFlag{
					Name:  "alpha",
					Usage: "Release for alpha version",
				},
			},
			ActionFunc: releaseDistribution,
		},
		{
			Name:  "mock",
			Usage: "Generate mock class",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-delete",
					Usage: "Generate mock class with delete previous generation",
				},
			},
			ActionFunc: generateMock,
		},
		{
			Name:       "readme",
			Usage:      "Generate readme document",
			ActionFunc: generateReadme,
		},
		{
			Name:       "docker",
			Usage:      "Docker utility",
			BeforeFunc: typienv.LoadEnvFile,
			SubCommands: []*typictx.Command{
				{
					Name:       "compose",
					Usage:      "Generate docker-compose.yaml",
					ActionFunc: dockerCompose,
				},
				{
					Name:  "up",
					Usage: "Create and start containers",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "no-compose",
							Usage: "Create and start containers without generate docker-compose.yaml",
						},
					},
					ActionFunc: dockerUp,
				},
				{
					Name:       "down",
					Usage:      "Stop and remove containers, networks, images, and volumes",
					ActionFunc: dockerDown,
				},
			},
		},
	}
	for _, module := range c.Modules {
		if module.Command != nil {
			cmds = append(cmds, module.Command)
		}
	}
	return
}

func buildBinary(ctx *typictx.ActionContext) error {
	binaryName := typienv.App.BinPath
	mainPackage := typienv.App.SrcPath
	return bash.GoBuild(binaryName, mainPackage)
}

func cleanProject(ctx *typictx.ActionContext) (err error) {
	if err = os.RemoveAll(typienv.Bin); err != nil {
		return
	}
	if err = os.RemoveAll(typienv.Metadata); err != nil {
		return
	}
	return filepath.Walk(typienv.Dependency.SrcPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return os.Remove(path)
		}
		return nil
	})
}

func runBinary(ctx *typictx.ActionContext) error {
	if !ctx.Cli.Bool("no-build") {
		buildBinary(ctx)
	}
	binaryPath := typienv.App.BinPath
	return bash.Run(binaryPath, []string(ctx.Cli.Args())...)
}

func runTesting(ctx *typictx.ActionContext) error {
	return bash.GoTest(ctx.TestTargets)
}

func generateMock(ctx *typictx.ActionContext) (err error) {
	err = bash.GoGet("github.com/golang/mock/mockgen")
	if err != nil {
		return
	}
	mockPkg := typienv.Mock
	if !ctx.Cli.Bool("no-delete") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}
	for _, mockTarget := range ctx.MockTargets {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]
		err = bash.RunGoBin("mockgen",
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
	return
}

func releaseDistribution(action *typictx.ActionContext) (err error) {
	if !action.Cli.Bool("no-test") {
		err = runTesting(action)
		if err != nil {
			return
		}
	}
	force := action.Cli.Bool("force")
	alpha := action.Cli.Bool("alpha")
	binaries, changeLogs, err := releaser.ReleaseDistribution(action.Release, force, alpha)
	if err != nil {
		return
	}
	if !action.Cli.Bool("no-github") {
		releaser.GithubRelease(binaries, changeLogs, action.Release, alpha)
	}
	return
}

func dockerCompose(ctx *typictx.ActionContext) (err error) {
	log.Info("Generate docker-compose.yml")
	dockerCompose := ctx.DockerCompose()
	d1, _ := yaml.Marshal(dockerCompose)
	return ioutil.WriteFile("docker-compose.yml", d1, 0644)
}

func dockerUp(ctx *typictx.ActionContext) (err error) {
	if !ctx.Cli.Bool("no-compose") {
		err = dockerCompose(ctx)
		if err != nil {
			return
		}
	}
	cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func dockerDown(ctx *typictx.ActionContext) (err error) {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// GenerateReadme for generate typical applical readme
func generateReadme(a *typictx.ActionContext) (err error) {
	readme0 := Readme{
		Title:       a.Name,
		Description: a.Description,
		Context:     a.Context,
	}
	log.Infof("Generate new %s", readmeFile)
	file, err := os.Create(readmeFile)
	if err != nil {
		return
	}
	defer file.Close()
	readme0.Markdown(file)
	return
}
