package typitask

import (
	"fmt"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"

	"github.com/google/go-github/v27/github"
)

type githubReleaser struct {
	typictx.Context
}

func (r *githubReleaser) IsReleased(service *github.RepositoriesService) bool {
	_, _, err := service.GetReleaseByTag(r,
		r.Release.Github.Owner,
		r.Release.Github.RepoName,
		r.ReleaseVersion())
	if err == nil {
		return true
	}
	return false
}

func (r *githubReleaser) CreateRelease(service *github.RepositoriesService, releaseNote string) (release *github.RepositoryRelease, err error) {
	release, _, err = service.CreateRelease(r,
		r.Release.Github.Owner,
		r.Release.Github.RepoName,
		&github.RepositoryRelease{
			Name:       github.String(fmt.Sprintf("%s - %s", r.Name, r.ReleaseVersion())),
			TagName:    github.String(r.ReleaseVersion()),
			Body:       github.String(releaseNote),
			Draft:      github.Bool(false),
			Prerelease: github.Bool(r.Release.Alpha),
		},
	)
	return
}

func (r *githubReleaser) Upload(service *github.RepositoriesService, repoID int64, binary string) (err error) {
	binaryPath := fmt.Sprintf("%s/%s", typienv.Release(), binary)

	var file *os.File
	file, err = os.Open(binaryPath)
	if err != nil {
		return
	}

	_, _, err = service.UploadReleaseAsset(r,
		r.Release.Github.Owner,
		r.Release.Github.RepoName,
		repoID,
		&github.UploadOptions{
			Name: binary,
		},
		file,
	)

	return
}
