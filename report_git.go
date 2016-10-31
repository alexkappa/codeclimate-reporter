package main

import (
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
	return fmt.Sprintf("Head: %s\nBranch: %s\nCommitted At: %s\n\n", g.Head, g.Branch, time.Unix(g.CommittedAt, 0).Format(time.RFC3339))
}

func collectGitInfo() (*Git, error) {
	cwd, _ := os.Getwd()

	repo, err := git.OpenRepository(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading git repository")
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading head")
	}
	// branch, err := ref.Name()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Failed reading branch name")
	// }
	commit, err := repo.LookupCommit(ref.Target())
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading commit")
	}

	it, err := repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil, errors.Wrap(err, "Failed iterator creationg")
	}

	for {
		branch, branchType, _ := it.Next()
		if branch == nil {
			break
		}
		name, err := branch.Name()
		if err != nil {
			break
		}
		fmt.Println(name, branchType)
	}

	for i := uint(0); i < commit.ParentCount(); i++ {
		fmt.Print("parent: ")
		fmt.Println(commit.Parent(i).TreeId())
		fmt.Print("parent id: ")
		fmt.Println(commit.ParentId(i))
	}
	fmt.Print("OWNER!!")
	owner := commit.Owner()
	refOwn, err := owner.Head()
	if err != nil {
		return nil, errors.Wrap(err, "Failed reading head")
	}
	fmt.Print("Try owner: ")
	fmt.Println(refOwn.Branch().Name())
	fmt.Println(repo)
	fmt.Print("Tree ID: ")
	fmt.Println(commit.TreeId())
	return &Git{
		Head:        ref.Target().String(),
		Branch:      ref.Shorthand(),
		CommittedAt: commit.Committer().When.Unix(),
	}, nil
}
