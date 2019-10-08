package releaser

import (
	"context"
	"fmt"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"

	"github.com/google/go-github/v27/github"
)

type githubReleaser struct {
	typictx.Release
	alpha bool
}

func (r *githubReleaser) IsReleased(ctx context.Context, service *github.RepositoriesService) bool {
	_, _, err := service.GetReleaseByTag(ctx,
		r.Release.Github.Owner,
		r.Release.Github.RepoName,
		r.ReleaseTag(r.alpha))
	if err == nil {
		return true
	}
	return false
}

func (r *githubReleaser) CreateRelease(ctx context.Context, service *github.RepositoriesService, releaseNote string) (release *github.RepositoryRelease, err error) {
	releaseTag := r.ReleaseTag(r.alpha)
	release, _, err = service.CreateRelease(ctx,
		r.Release.Github.Owner,
		r.Release.Github.RepoName,
		&github.RepositoryRelease{
			Name:       github.String(fmt.Sprintf("%s - %s", r.ReleaseName(), releaseTag)),
			TagName:    github.String(releaseTag),
			Body:       github.String(releaseNote),
			Draft:      github.Bool(false),
			Prerelease: github.Bool(r.alpha),
		},
	)
	return
}

func (r *githubReleaser) Upload(ctx context.Context, service *github.RepositoriesService, repoID int64, binary string) (err error) {
	binaryPath := fmt.Sprintf("%s/%s", typienv.Release, binary)
	var file *os.File
	file, err = os.Open(binaryPath)
	if err != nil {
		return
	}
	_, _, err = service.UploadReleaseAsset(ctx,
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
