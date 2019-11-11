package typictx

import (
	"errors"
	"fmt"
	"strings"
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
