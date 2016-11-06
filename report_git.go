package main

import (
	"bytes"
	"fmt"
	"os"
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

	bch, err := ref.Branch().Name()
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading branch name")
	}

	cmt, err := rep.LookupCommit(ref.Target())
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading commit")
	}

	return &Git{
		Branch:      bch,
		Head:        ref.Target().String(),
		CommittedAt: cmt.Committer().When.Unix(),
	}, nil
}
