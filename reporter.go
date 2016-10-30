package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	CodeClimateRepoToken = os.Getenv("CODECLIMATE_REPO_TOKEN")
	CodeClimateAPIHost   = os.Getenv("CODECLIMATE_API_HOST")
)

func init() {
	if CodeClimateAPIHost == "" {
		CodeClimateAPIHost = "codeclimate.com"
	}
}

type reporter struct {
	*http.Client
}

func newReporter(skipTLSVerify bool) *reporter {
	return &reporter{
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify},
			},
		},
	}
}

func (rep *reporter) send(r *Report) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(r)
	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", "https://"+CodeClimateAPIHost+"/test_reports", &body)
	request.Header.Set("User-Agent", "Code Climate (Go Test Reporter "+Version+")")
	request.Header.Set("Content-Type", "application/json")

	response, err := rep.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == 401 {
		return fmt.Errorf("an invalid CODECLIMATE_REPO_TOKEN token was specified")
	}

	return nil
}

type Report struct {
	Partial     bool         `json:"partial"`
	RunAt       int64        `json:"run_at"`
	RepoToken   string       `json:"repo_token"`
	Environment *Env         `json:"environment"`
	Git         *Git         `json:"git"`
	CI          *CI          `json:"ci_service"`
	SourceFiles []SourceFile `json:"source_files"`
}

func makeReport(r io.Reader) (*Report, error) {
	env, err := collectEnv()
	if err != nil {
		return nil, err
	}

	git, err := collectGitInfo()
	if err != nil {
		return nil, err
	}

	ci, err := collectCIInfo()
	if err != nil {
		return nil, err
	}

	cov, err := collectCoverage(r)
	if err != nil {
		return nil, err
	}

	src, err := collectSource(cov)
	if err != nil {
		return nil, err
	}

	return &Report{
		Partial:     false,
		RunAt:       time.Now().Unix(),
		RepoToken:   CodeClimateRepoToken,
		Environment: env,
		Git:         git,
		CI:          ci,
		SourceFiles: src,
	}, nil
}
