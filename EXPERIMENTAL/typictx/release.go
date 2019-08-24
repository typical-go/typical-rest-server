package typictx

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/git"
)

// Release setting
type Release struct {
	Versioning
	Name    string
	Version string
	Alpha   bool
	GoOS    []string
	GoArch  []string
	Github  *Github
}

// Github contain github information
type Github struct {
	Owner    string
	RepoName string
}

// Versioning setting
type Versioning struct {
	WithGitBranch       bool
	WithLatestGitCommit bool
}

// Tag to get release tag
func (r *Release) Tag() string {
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

// ReleaseBinary to get release binary
func (r *Release) ReleaseBinary(os1, arch string) string {
	name := r.Name
	if name == "" {
		dir, _ := os.Getwd()
		name = filepath.Base(dir)
	}
	return strings.Join([]string{name, r.Tag(), os1, arch}, "_")
}
