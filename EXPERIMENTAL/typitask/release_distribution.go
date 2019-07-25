package typitask

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"gopkg.in/src-d/go-git.v4"

	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v27/github"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// ReleaseDistribution to release binary distribution
func ReleaseDistribution(ctx typictx.ActionContext) (err error) {
	gitRepo, err := git.PlainOpen(".")
	if err != nil {
		return
	}

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
		log.Info("Generate default note from git logs.")
		note = defaultNote(gitRepo)
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

func defaultNote(gitRepo *git.Repository) string {
	var builder strings.Builder

	tagrefs, _ := gitRepo.Tags()
	latestTag, _ := tagrefs.Next()
	tagrefs.Close()

	gitLogs, _ := gitRepo.Log(&git.LogOptions{})
	defer gitLogs.Close()
	for {
		commit, err := gitLogs.Next()
		if err != nil {
			break
		}

		if commit.Hash == latestTag.Hash() {
			break
		}

		shortHash := commit.Hash.String()[0:8]
		message := strings.TrimSpace(commit.Message)

		if !ignoredMessage(message) {
			builder.WriteString(fmt.Sprintf("%s %s\n", shortHash, message))
		}
	}

	return builder.String()
}

func ignoredMessage(message string) bool {
	// TODO: ignore everything start with small-case character
	return strings.HasPrefix(message, "Merge")
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
	repo := githubDetail.RepoName

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, releaseInfo.Version)
	if err == nil {
		if force {
			log.Infof("Force release detected; Delete existing release for %s/%s (%s)", owner, repo, releaseInfo.Version)
			_, err = client.Repositories.DeleteRelease(ctx, owner, repo, *release.ID)
			if err != nil {
				return
			}
		} else {
			log.Infof("Release for %s/%s (%s) already exist", owner, repo, releaseInfo.Version)
			return nil
		}
	}

	log.Infof("Create github release for %s/%s", owner, repo)
	release, _, err = client.Repositories.CreateRelease(ctx, owner, repo, releaseInfo.Data())
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
			repo,
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
