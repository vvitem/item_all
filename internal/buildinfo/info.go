package buildinfo

import "runtime"

const Name = "ItemAll"

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

func Current() Info {
	return Info{
		Name:      Name,
		Version:   valueOrDefault(Version, "dev"),
		Commit:    valueOrDefault(Commit, "unknown"),
		BuildTime: valueOrDefault(BuildTime, "unknown"),
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func valueOrDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
