package typictx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/git"
)

// Release setting
type Release struct {
	Tagging
	Name    string
	Version string
	Targets []string
	Github  *Github
}

// Github contain github information
type Github struct {
	Owner    string
	RepoName string
}

// Tagging setting
type Tagging struct {
	WithGitBranch       bool
	WithLatestGitCommit bool
}

// Validate the release
func (r *Release) Validate() (err error) {
	if len(r.Targets) < 1 {
		return errors.New("Missing 'Targets'")
	}
	for _, target := range r.Targets {
		if !strings.Contains(target, "/") {
			return fmt.Errorf("Missing '/' in target '%s'", target)
		}
	}
	return
}

// ReleaseTag to get release tag
func (r *Release) ReleaseTag(alpha bool) string {
	var builder strings.Builder
	builder.WriteString("v")
	builder.WriteString(r.Version)
	if alpha {
		builder.WriteString("-alpha")
	}
	if r.WithGitBranch {
		builder.WriteString("_")
		builder.WriteString(git.Branch())
	}
	if r.WithLatestGitCommit {
		builder.WriteString("_")
		builder.WriteString(git.LatestCommit())
	}
	return builder.String()
}

// ReleaseName to get release name
func (r *Release) ReleaseName() string {
	name := r.Name
	if name == "" {
		dir, _ := os.Getwd()
		name = filepath.Base(dir)
	}
	return name
}

// ReleaseBinary to get release binary
func (r *Release) ReleaseBinary(os1, arch string, alpha bool) string {
	return strings.Join([]string{
		r.ReleaseName(),
		r.ReleaseTag(alpha),
		os1,
		arch,
	}, "_")
}
