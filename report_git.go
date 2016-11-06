package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alexkappa/errors"
	git "github.com/libgit2/git2go"
)

type Git struct {
	Head        string `json:"head"`
	Branch      string `json:"branch"`
	CommittedAt int64  `json:"committed_at"`
}

func (g *Git) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Head: %s\n", g.Head)
	fmt.Fprintf(&buf, "Branch: %s\n", g.Branch)
	fmt.Fprintf(&buf, "Committed At: %s\n\n", time.Unix(g.CommittedAt, 0).Format(time.RFC3339))
	return buf.String()
}

func collectGitInfo() (*Git, error) {

	cwd, _ := os.Getwd()
	rep, err := git.OpenRepository(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading git repository")
	}

	ref, err := rep.Head()
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading head")
	}
	commit, err := rep.LookupCommit(ref.Target())
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading commit")
	}

	bch, err := collectGitBranch(fmt.Sprintf("%s", commit.Id()))
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading branch name")
	}

	cmt, err := rep.LookupCommit(ref.Target())
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading commit")
	}

	return &Git{
		Head:        ref.Target().String(),
		Branch:      bch,
		CommittedAt: cmt.Committer().When.Unix(),
	}, nil
}

func collectGitBranch(commitID string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("git", "branch", "-r", "--contains", commitID)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	}
	s := bufio.NewScanner(&stdout)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "*") {
			return strings.TrimPrefix(s.Text(), "* "), nil
		}
	}
	if s.Err() != nil {
		return "", s.Err()
	}
	return "", errors.New("Failed parsing `git branch` output")
}
