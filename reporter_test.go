package main

import (
	"testing"
	"time"
)

func TestReport(t *testing.T) {
	_ = Report{
		Partial:   false,
		RunAt:     time.Now().Unix(),
		RepoToken: "abc",
	}
}
