package typictx

import (
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
	Alpha   bool
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

// ReleaseTag to get release tag
func (r *Release) ReleaseTag() string {
	var builder strings.Builder
	builder.WriteString("v")
	builder.WriteString(r.Version)
	if r.Alpha {
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
func (r *Release) ReleaseBinary(os1, arch string) string {

	return strings.Join([]string{r.ReleaseName(), r.ReleaseTag(), os1, arch}, "_")
}
