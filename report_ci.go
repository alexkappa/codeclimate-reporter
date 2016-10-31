package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type CI struct {
	Branch      string `json:"branch"`
	BuildID     string `json:"build_identifier"`
	BuildURL    string `json:"build_url"`
	CommitSHA   string `json:"commit_sha"`
	Name        string `json:"name"`
	PullRequest string `json:"pull_request"`
	WorkerID    string `json:"worker_id"`
}

func (c *CI) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Name: %s\n", c.Name)
	fmt.Fprintf(&buf, "BuildID: %s\n", c.BuildID)
	fmt.Fprintf(&buf, "BuildURL: %s\n", c.BuildURL)
	fmt.Fprintf(&buf, "CommitSHA: %s\n", c.CommitSHA)
	fmt.Fprintf(&buf, "PullRequest: %s\n", c.PullRequest)
	fmt.Fprintf(&buf, "WorkerID: %s\n\n", c.WorkerID)
	return buf.String()
}

func collectCIInfo() (*CI, error) {
	ci := &CI{}

	switch {
	case os.Getenv("TRAVIS") != "":
		ci.Name = "travis-ci"
		ci.Branch = os.Getenv("TRAVIS_BRANCH")
		ci.BuildID = os.Getenv("TRAVIS_JOB_ID")
		ci.PullRequest = os.Getenv("TRAVIS_PULL_REQUEST")
	case os.Getenv("CIRCLECI") != "":
		ci.Name = "circleci"
		ci.BuildID = os.Getenv("CIRCLE_BUILD_NUM")
		ci.Branch = os.Getenv("CIRCLE_BRANCH")
		ci.CommitSHA = os.Getenv("CIRCLE_SHA1")
	case os.Getenv("SEMAPHORE") != "":
		ci.Name = "semaphore"
		ci.Branch = os.Getenv("BRANCH_NAME")
		ci.BuildID = os.Getenv("SEMAPHORE_BUILD_NUMBER")
	case os.Getenv("JENKINS_URL") != "":
		ci.Name = "jenkins"
		ci.BuildID = os.Getenv("BUILD_NUMBER")
		ci.BuildURL = os.Getenv("BUILD_URL")
		ci.Branch = os.Getenv("GIT_BRANCH")
		ci.CommitSHA = os.Getenv("GIT_COMMIT")
	case os.Getenv("TDDIUM") != "":
		ci.Name = "tddium"
		ci.BuildID = os.Getenv("TDDIUM_SESSION_ID")
		ci.WorkerID = os.Getenv("TDDIUM_TID")
	case os.Getenv("WERCKER") != "":
		ci.Name = "wercker"
		ci.BuildID = os.Getenv("WERCKER_BUILD_ID")
		ci.BuildURL = os.Getenv("WERCKER_BUILD_URL")
		ci.Branch = os.Getenv("WERCKER_GIT_BRANCH")
		ci.CommitSHA = os.Getenv("WERCKER_GIT_COMMIT")
	case strings.Contains(os.Getenv("CI_NAME"), "codeship"):
		ci.Name = "codeship"
		ci.BuildID = os.Getenv("CI_BUILD_NUMBER")
		ci.BuildURL = os.Getenv("CI_BUILD_URL")
		ci.Branch = os.Getenv("CI_BRANCH")
		ci.CommitSHA = os.Getenv("CI_COMMIT_ID")
	case os.Getenv("APPVEYOR") != "":
		ci.Name = "appveyor"
		ci.BuildID = os.Getenv("APPVEYOR_BUILD_NUMBER")
		ci.Branch = os.Getenv("APPVEYOR_REPO_BRANCH")
		ci.CommitSHA = os.Getenv("APPVEYOR_REPO_COMMIT")
		ci.PullRequest = os.Getenv("APPVEYOR_PULL_REQUEST_NUMBER")
	case os.Getenv("BUILDKITE") != "":
		ci.Name = "buildkite"
		ci.BuildID = os.Getenv("BUILDKITE_BUILD_ID")
		ci.BuildURL = os.Getenv("BUILDKITE_BUILD_URL")
		ci.Branch = os.Getenv("BUILDKITE_BRANCH")
		ci.CommitSHA = os.Getenv("BUILDKITE_COMMIT")
	}

	return ci, nil
}
