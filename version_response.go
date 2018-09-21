package main

import "runtime"

type VersionResponse struct {
	Description string `json:"description"`
	GitCommit   string `json:"git_commit"`
	GoVersion   string `json:"go_version"`
	Name        string `json:"name"`
	OSArch      string `json:"os_arch"`
	Source      string `json:"source"`
}

func newVersionResponse() VersionResponse {
	return VersionResponse{
		Description: description,
		GitCommit:   gitCommit,
		GoVersion:   runtime.Version(),
		Name:        name,
		OSArch:      runtime.GOOS + "/" + runtime.GOARCH,
		Source:      source,
	}
}
