package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Status is same with `git status --porcelain`
func Status(files ...string) string {
	args := []string{"status"}
	args = append(args, files...)
	args = append(args, "--porcelain")
	status, err := Git(args...)
	if err != nil {
		return err.Error()
	}
	return status
}

// Fetch is same with `get fetch`
func Fetch() error {
	return exec.Command("git", "fetch").Run()
}

// LatestTag to get latest tag and its hash key
func LatestTag() string {
	tag, err := Git("describe", "--tags", "--abbrev=0")
	if err != nil {
		return ""
	}
	return tag
}

// Logs of commits
func Logs(from string) []string {
	data, err := Git("--no-pager", "log", fmt.Sprintf("%s..HEAD", from), "--oneline")
	if err != nil {
		return []string{}
	}
	return strings.Split(data, "\n")
}

// Push files to git repo
func Push(commitMessage string, files ...string) (err error) {
	args := []string{"add"}
	args = append(args, files...)
	_, err = Git(args...)
	if err != nil {
		return
	}
	_, err = Git("commit", "-m", commitMessage)
	if err != nil {
		return
	}
	_, err = Git("push")
	return
}

// Branch to return current branch
func Branch() string {
	branch, err := Git("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}
	return branch
}

// LatestCommit return latest commit in short hash
func LatestCommit() string {
	commit, err := Git("rev-parse", "--short", "HEAD")
	if err != nil {
		return ""
	}
	return commit
}

// Git to execution git command
func Git(args ...string) (s string, err error) {
	var builder strings.Builder
	cmd := exec.Command("git", args...)
	cmd.Stdout = &builder
	cmd.Stderr = &builder
	err = cmd.Run()
	s = strings.TrimSuffix(builder.String(), "\n")
	return
}
