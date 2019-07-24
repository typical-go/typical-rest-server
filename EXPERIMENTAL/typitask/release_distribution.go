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
	// RunTest(ctx)
	// GenerateReadme(ctx)

	goos := []string{"linux", "darwin"}
	goarch := []string{"amd64"}
	mainPackage := typienv.AppMainPackage()

	var binaries []string

	for _, os1 := range goos {
		for _, arch := range goarch {
			// TODO: using ldflags
			binary := fmt.Sprintf("%s_%s_%s_%s",
				ctx.Typical.BinaryNameOrDefault(),
				ctx.Typical.Version,
				os1,
				arch)

			binaryPath := fmt.Sprintf("%s/%s", typienv.Release(), binary)

			log.Infof("Create release for %s/%s: %s", os1, arch, binary)
			os.Setenv("GOOS", os1)
			os.Setenv("GOARCH", arch)
			bash.GoBuild(binaryPath, mainPackage)

			binaries = append(binaries, binary)
		}
	}

	githubKey := ctx.Cli.String("github-token")
	if githubKey != "" {
		err = releaseToGithub(githubKey, binaries)
	}

	return nil
}

func releaseToGithub(token string, binaries []string) (err error) {
	log.Info("Release to Github")

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// TOOD: get from typical context
	title := "Typical Go Server - v0.0.1-alpha"
	currentTag := "v0.0.1-alpha"
	preRelease := true
	owner := "typical-go"
	name := "typical-rest-server"

	// TODO: list of commit message
	body := ""

	var release *github.RepositoryRelease

	var data = &github.RepositoryRelease{
		Name:       github.String(title),
		TagName:    github.String(currentTag),
		Body:       github.String(body),
		Draft:      github.Bool(false),
		Prerelease: github.Bool(preRelease),
	}

	release, _, err = client.Repositories.CreateRelease(ctx, owner, name, data)

	if err != nil {
		return
	}

	for _, binary := range binaries {
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
