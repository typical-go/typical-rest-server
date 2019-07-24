package typitask

import (
	"context"
	"fmt"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"

	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v27/github"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// ReleaseDistribution to release binary distribution
func ReleaseDistribution(ctx typictx.ActionContext) (err error) {
	if !ctx.Cli.Bool("no-test") {
		RunTest(ctx)
	}

	if !ctx.Cli.Bool("no-readme") {
		GenerateReadme(ctx)
		// TODO: check git diff and commit if readme is updated
	}

	goos := []string{"linux", "darwin"}
	goarch := []string{"amd64"}
	mainPackage := typienv.AppMainPackage()

	alpha := ctx.Cli.Bool("alpha")
	version := fmt.Sprintf("v%s", ctx.Typical.Version)
	if alpha {
		version = fmt.Sprintf("%s-alpha", version)
	}

	var binaries []string

	for _, os1 := range goos {
		for _, arch := range goarch {
			// TODO: consider to using ldflags
			binary := fmt.Sprintf("%s_%s_%s_%s",
				ctx.Typical.BinaryNameOrDefault(),
				version,
				os1,
				arch)

			binaryPath := fmt.Sprintf("%s/%s", typienv.Release(), binary)

			log.Infof("Create release binary for %s/%s at %s", os1, arch, binary)
			os.Setenv("GOOS", os1)
			os.Setenv("GOARCH", arch)
			bash.GoBuild(binaryPath, mainPackage)

			binaries = append(binaries, binary)
		}
	}

	note := ctx.Cli.String("note")
	if note == "" {
		// TODO: generate default note from git log
	}

	githubKey := ctx.Cli.String("github-token")
	if githubKey != "" {
		err = releaseToGithub(
			ctx.Typical.Github,
			githubKey,
			githubReleaseInfo{
				ApplicationName: ctx.Typical.Name,
				Binaries:        binaries,
				Version:         version,
				Alpha:           alpha,
				Note:            note,
			},
			ctx.Cli.Bool("force"),
		)
	}

	return
}

func releaseToGithub(githubDetail *typictx.Github, token string, releaseInfo githubReleaseInfo, force bool) (err error) {
	if githubDetail == nil {
		return fmt.Errorf("Missing Github in typical context")
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner := githubDetail.Owner
	name := githubDetail.Name

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, name, releaseInfo.Version)
	if err == nil {
		if force {
			log.Infof("Force release detected; Delete existing release for %s/%s (%s)", owner, name, releaseInfo.Version)
			_, err = client.Repositories.DeleteRelease(ctx, owner, name, *release.ID)
			if err != nil {
				return
			}
		} else {
			log.Infof("Release for %s/%s (%s) already exist", owner, name, releaseInfo.Version)
			return nil
		}
	}

	log.Infof("Create github release for %s/%s", owner, name)
	release, _, err = client.Repositories.CreateRelease(ctx, owner, name, releaseInfo.Data())
	if err != nil {
		return
	}

	for _, binary := range releaseInfo.Binaries {
		var file *os.File
		binaryPath := fmt.Sprintf("%s/%s", typienv.Release(), binary)
		log.Info("Upload release asset: " + binaryPath)

		file, err = os.Open(binaryPath)
		if err != nil {
			return
		}

		_, _, err = client.Repositories.UploadReleaseAsset(
			ctx,
			owner,
			name,
			release.GetID(),
			&github.UploadOptions{
				Name: binary,
			},
			file,
		)
	}

	return
}

type githubReleaseInfo struct {
	ApplicationName string
	Binaries        []string
	Version         string
	Alpha           bool
	Note            string
}

func (i *githubReleaseInfo) Data() *github.RepositoryRelease {
	return &github.RepositoryRelease{
		Name:       github.String(fmt.Sprintf("%s - %s", i.ApplicationName, i.Version)),
		TagName:    github.String(i.Version),
		Body:       github.String(i.Note),
		Draft:      github.Bool(false),
		Prerelease: github.Bool(i.Alpha),
	}
}
