package main

import (
	"fmt"
	"os"
)

type Env struct {
	Pwd            string `json:"pwd"`
	PackageVersion string `json:"package_version"`
}

func (e *Env) String() string {
	return fmt.Sprintf("Pwd: %s\nPackage Version: %s\n\n", e.Pwd, e.PackageVersion)
}

func collectEnv() (*Env, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Env{
		Pwd:            pwd,
		PackageVersion: Version,
	}, nil
}
