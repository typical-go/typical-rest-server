package typirelease

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v27/github"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/git"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"golang.org/x/oauth2"
)

// Release the distribution
func Release(ctx *typictx.Context) (err error) {
	git.Fetch()
	defer git.Fetch()
	releaseVersion := releaseVersion(ctx)
	latestTag := git.LatestTag()
	if latestTag == releaseVersion {
		log.Infof("%s already released", latestTag)
		return nil
	}
	changeLogs := git.Logs(latestTag)
	if len(changeLogs) < 1 {
		log.Info("No change to be released")
		return nil
	}
	for _, changeLog := range changeLogs {
		log.Infof("Change Log: %s", changeLog)
	}
	binaries, err := buildReleaseBinaries(ctx, releaseVersion)
	if err != nil {
		return
	}
	if ctx.Release.Github != nil {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			return errors.New("Environment 'GITHUB_TOKEN' is missing")
		}

		owner := ctx.Release.Github.Owner
		repo := ctx.Release.Github.RepoName

		ctx0 := context.Background()
		client := github.NewClient(oauth2.NewClient(
			ctx0,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		))

		releaser := githubReleaser{ctx}
		if releaser.IsReleased(ctx0, client.Repositories) {
			log.Infof("Release for %s/%s (%s) already exist", owner, repo, releaseVersion)
			return
		}

		log.Info("Generate release note")
		var releaseNote strings.Builder
		for _, changelog := range changeLogs {
			if !ignoring(changelog) {
				releaseNote.WriteString(changelog)
				releaseNote.WriteString("\n")
			}
		}

		log.Infof("Create github release for %s/%s", owner, repo)
		var release *github.RepositoryRelease
		release, err = releaser.CreateRelease(ctx0, client.Repositories, releaseNote.String())
		if err != nil {
			return
		}

		for _, binary := range binaries {
			log.Infof("Upload asset: %s", binary)
			err = releaser.Upload(ctx0, client.Repositories, *release.ID, binary)
			if err != nil {
				return
			}
		}
	}

	return
}

func buildReleaseBinaries(ctx *typictx.Context, version string) (binaries []string, err error) {
	if len(ctx.GoOS) < 0 {
		err = errors.New("Missing 'GoOS' in Typical Context")
		return
	}

	if len(ctx.GoArch) < 0 {
		err = errors.New("Missing 'GoArch' in Typical Context")
		return
	}

	mainPackage := typienv.AppMainPackage()
	for _, os1 := range ctx.GoOS {
		for _, arch := range ctx.GoArch {
			// TODO: consider to using ldflags
			binary := fmt.Sprintf("%s_%s_%s_%s",
				ctx.BinaryNameOrDefault(), version, os1, arch)

			binaryPath := fmt.Sprintf("%s/%s", typienv.Release(), binary)

			log.Infof("Create release binary for %s/%s at %s", os1, arch, binary)
			os.Setenv("GOOS", os1)
			os.Setenv("GOARCH", arch)
			err = bash.GoBuild(binaryPath, mainPackage, "-w", "-s")
			if err != nil {
				return
			}

			binaries = append(binaries, binary)
		}
	}
	return
}

func cleanMessage(message string) string {
	iCoAuthor := strings.Index(message, "Co-Authored-By")
	if iCoAuthor > 0 {
		message = message[0:strings.Index(message, "Co-Authored-By")]
	}
	message = strings.TrimSpace(message)
	return message
}

func ignoring(changelog string) bool {
	message := cleanMessage(changelog[7:])
	lowerMessage := strings.ToLower(message)
	return strings.HasPrefix(lowerMessage, "merge") ||
		strings.HasPrefix(lowerMessage, "bump") ||
		strings.HasPrefix(lowerMessage, "revision") ||
		strings.HasPrefix(lowerMessage, "generate") ||
		strings.HasPrefix(lowerMessage, "wip")
}

func releaseVersion(c *typictx.Context) (version string) {
	version = fmt.Sprintf("v%s", c.Version)
	if c.Release.Alpha {
		version = fmt.Sprintf("%s-alpha", version)
	}
	return
}
