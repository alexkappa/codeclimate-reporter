package main

import "os"

type Env struct {
	Pwd            string `json:"pwd"`
	PackageVersion string `json:"package_version"`
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
