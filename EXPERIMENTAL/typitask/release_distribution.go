package typitask

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v27/github"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// ReleaseDistribution to release binary distribution
func ReleaseDistribution(ctx *typictx.ActionContext) (err error) {
	// NOTE: git fetch in beginning and after to make local is up2date
	exec.Command("git", "fetch").Run()
	defer exec.Command("git", "fetch").Run()

	goos := []string{"linux", "darwin"}
	goarch := []string{"amd64"}
	mainPackage := typienv.AppMainPackage()

	alpha := ctx.Cli.Bool("alpha")
	version := fmt.Sprintf("v%s", ctx.Typical.Version)
	if alpha {
		version = fmt.Sprintf("%s-alpha", version)
	}

	gitRepo, err := git.PlainOpen(".")
	if err != nil {
		return
	}

	worktree, err := gitRepo.Worktree()
	if err != nil {
		return
	}

	status, err := worktree.Status()
	if err != nil {
		return
	}

	latestTag := latestTag(gitRepo)
	if latestTag != nil && latestTag.Name().Short() == version {
		log.Infof("%s already released", version)
		return nil
	}

	if !status.IsClean() {
		log.Info("Please submit uncommitted change first")
		return nil
	}

	changes := changeLogs(gitRepo, latestTag)
	if len(changes) < 1 {
		log.Info("No change to be released")
		return nil
	}

	for _, change := range changes {
		log.Infof("Change Log: %s", change)
	}

	if !ctx.Cli.Bool("no-test") {
		err = RunTest(ctx)
		if err != nil {
			return
		}
	}

	if !ctx.Cli.Bool("no-readme") {
		err = GenerateReadme(ctx)
		if err != nil {
			return
		}
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
			err = bash.GoBuild(binaryPath, mainPackage, "-w", "-s")
			if err != nil {
				return
			}

			binaries = append(binaries, binary)
		}
	}

	note := ctx.Cli.String("note")
	if note == "" {
		log.Info("Using git change log as release note.")
		note = strings.Join(changes, "")
	}

	githubKey := ctx.Cli.String("github-token")
	if githubKey != "" {
		err = releaseToGithub(
			ctx,
			githubKey,
			githubReleaseInfo{
				ApplicationName: ctx.Typical.Name,
				Binaries:        binaries,
				Version:         version,
				Alpha:           alpha,
				Note:            note,
			},
		)
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

// TODO: change logs should be unfiltered. Filter only when generate release page
func changeLogs(gitRepo *git.Repository, latestTag *plumbing.Reference) (changes []string) {

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

		if !ignoredMessage(message) {
			change := fmt.Sprintf("%s %s\n", shortHash, message)
			changes = append(changes, change)
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

func ignoredMessage(message string) bool {
	lowerMessage := strings.ToLower(message)

	// TODO: consider to ignore message if start with small-case character

	return strings.HasPrefix(lowerMessage, "merge") ||
		strings.HasPrefix(lowerMessage, "bump") ||
		strings.HasPrefix(lowerMessage, "revision") ||
		strings.HasPrefix(lowerMessage, "generate") ||
		strings.HasPrefix(lowerMessage, "wip")
}

func releaseToGithub(ctx *typictx.ActionContext, token string, releaseInfo githubReleaseInfo) (err error) {
	githubDetail := ctx.Typical.Github
	if githubDetail == nil {
		return fmt.Errorf("Missing Github in typical context")
	}

	client := github.NewClient(oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	))

	owner := githubDetail.Owner
	repo := githubDetail.RepoName

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, releaseInfo.Version)
	if err == nil {
		if ctx.Cli.Bool("force") {
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
