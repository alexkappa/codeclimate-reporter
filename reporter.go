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

func (rep *reporter) send(r *Report) (*http.Request, *http.Response, error) {
	var body bytes.Buffer

	enc := json.NewEncoder(&body)
	enc.SetIndent("", "  ")
	err := enc.Encode(r)
	if err != nil {
		return nil, nil, err
	}

	req, _ := http.NewRequest("POST", "https://"+CodeClimateAPIHost+"/test_reports", &body)
	req.Header.Set("User-Agent", "Code Climate (Go Test Reporter "+Version+")")
	req.Header.Set("Content-Type", "application/json")

	res, err := rep.Do(req)
	if err != nil {
		return req, res, err
	}

	if res.StatusCode == 401 {
		return req, res, fmt.Errorf("an invalid CODECLIMATE_REPO_TOKEN token was specified")
	}

	return req, res, nil
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

func (r *Report) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Code Coverage Report\n")
	fmt.Fprintf(&buf, "====================\n\n")
	fmt.Fprintf(&buf, "Run At: %s\n", time.Unix(r.RunAt, 0).Format(time.RFC3339))
	fmt.Fprintf(&buf, "Repo Token: %s\n\n", CodeClimateRepoToken)
	fmt.Fprintf(&buf, "Environment\n")
	fmt.Fprintf(&buf, "-----------\n\n")
	buf.WriteString(r.Environment.String())
	fmt.Fprintf(&buf, "Git\n")
	fmt.Fprintf(&buf, "---\n\n")
	buf.WriteString(r.Git.String())
	fmt.Fprintf(&buf, "Continuous Integration\n")
	fmt.Fprintf(&buf, "----------------------\n\n")
	buf.WriteString(r.CI.String())
	fmt.Fprintf(&buf, "Source Files\n")
	fmt.Fprintf(&buf, "------------\n\n")
	for _, file := range r.SourceFiles {
		buf.WriteString(file.String())
	}
	return buf.String()
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
