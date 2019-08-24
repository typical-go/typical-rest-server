package typirelease

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v27/github"
	"github.com/labstack/gommon/log"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Release the distribution
func Release(ctx *typictx.Context) (err error) {
	// NOTE: git fetch in beginning and after to make local is up2date
	exec.Command("git", "fetch").Run()
	defer exec.Command("git", "fetch").Run()

	gitRepo, err := git.PlainOpen(".")
	if err != nil {
		return
	}

	latestTag := latestTag(gitRepo)
	if latestTag != nil && latestTag.Name().Short() == releaseVersion(ctx) {
		log.Infof("%s already released", releaseVersion(ctx))
		return nil
	}

	worktree, _ := gitRepo.Worktree()
	status, _ := worktree.Status()
	if !status.IsClean() {
		return fmt.Errorf("Please submit uncommitted change first:\n%s", status.String())
	}

	changelogs := listChangeLogs(gitRepo, latestTag)
	if len(changelogs) < 1 {
		log.Info("No change to be released")
		return nil
	}

	for _, changelog := range changelogs {
		log.Infof("Change Log: %s", changelog)
	}

	binaries, err := buildReleaseBinaries(ctx)
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
			log.Infof("Release for %s/%s (%s) already exist", owner, repo, releaseVersion(ctx))
			return
		}

		log.Info("Generate release note")
		var releaseNote strings.Builder
		for _, changelog := range changelogs {
			if !ignoring(changelog.Message) {
				releaseNote.WriteString(changelog.String())
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

func buildReleaseBinaries(ctx *typictx.Context) (binaries []string, err error) {
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
				ctx.BinaryNameOrDefault(),
				releaseVersion(ctx),
				os1,
				arch)

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

func latestTag(gitRepo *git.Repository) (latestTag *plumbing.Reference) {
	tagrefs, _ := gitRepo.Tags()
	defer tagrefs.Close()
	tagrefs.ForEach(func(tagRef *plumbing.Reference) error {
		latestTag = tagRef
		return nil
	})

	return
}

func listChangeLogs(gitRepo *git.Repository, latestTag *plumbing.Reference) (logs []changelog) {
	gitLogs, _ := gitRepo.Log(&git.LogOptions{})
	defer gitLogs.Close()
	for {
		commit, err := gitLogs.Next()
		if err != nil {
			break
		}
		if latestTag != nil && commit.Hash == latestTag.Hash() {
			break
		}
		shortHash := commit.Hash.String()[0:8]
		message := cleanMessage(commit.Message)

		logs = append(logs, changelog{
			Hash:    shortHash,
			Message: message,
		})
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

func ignoring(message string) bool {
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
