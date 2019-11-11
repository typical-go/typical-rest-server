package releaser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v27/github"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/git"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"golang.org/x/oauth2"
)

// Releaser responsible to release distruction
type Releaser struct {
	typictx.Release
	Force bool
	Alpha bool
}

// Git release and return change logs
func (r *Releaser) Git() (changeLogs []string, err error) {
	if len(r.Targets) < 0 {
		msg := "Missing 'Targets' in Typical Context; The format should be '$GOOS/$GOARCH'"
		log.Error(msg)
		err = errors.New(msg)
		return
	}
	git.Fetch()
	defer git.Fetch()
	status := git.Status()
	if !r.Force && status != "" {
		err = fmt.Errorf("Please commit changes first:\n%s", status)
		return
	}
	latestTag := git.LatestTag()
	if !r.Force && latestTag == r.ReleaseTag(r.Alpha) {
		log.Errorf("%s already released", latestTag)
		return
	}
	changeLogs = git.Logs(latestTag)
	if !r.Force && len(changeLogs) < 1 {
		msg := "No change to be released"
		log.Errorf(msg)
		err = errors.New(msg)
		return
	}
	return
}

// Distribution to release the distribution
func (r *Releaser) Distribution() (binaries []string, err error) {
	mainPackage := typienv.App.SrcPath
	for _, target := range r.Targets {
		chunks := strings.Split(target, "/")
		if len(chunks) != 2 {
			err = fmt.Errorf("Invalid target '%s': it should be '$GOOS/$GOARCH'", target)
			return
		}
		binary := r.ReleaseBinary(chunks[0], chunks[1], r.Alpha)
		binaryPath := fmt.Sprintf("%s/%s", typienv.Release, binary)
		log.Infof("Create release binary for %s: %s", target, binaryPath)
		// TODO: support cgo
		envs := []string{
			"GOOS=" + chunks[0],
			"GOARCH=" + chunks[1],
		}
		if err = bash.GoBuild(binaryPath, mainPackage, envs...); err != nil {
			return
		}
		binaries = append(binaries, binary)
	}
	return
}

// GithubRelease for github release
func GithubRelease(binaries, changeLogs []string, rel typictx.Release, alpha bool) (err error) {
	if rel.Github == nil {
		return
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Environment 'GITHUB_TOKEN' is missing")
	}
	owner := rel.Github.Owner
	repo := rel.Github.RepoName
	ctx0 := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx0, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	releaser := githubReleaser{rel, alpha}
	if releaser.IsReleased(ctx0, client.Repositories) {
		log.Infof("Release for %s/%s (%s) already exist", owner, repo, rel.ReleaseTag(alpha))
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
	return
}
