package main

import (
	"fmt"

	"github.com/vvitem/item_all/internal/buildinfo"
)

const tagline = "Local-first SafeOps"
const appStatusReady = "ready"

type App struct{}

type AppInfo struct {
	Name      string `json:"name"`
	Tagline   string `json:"tagline"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	Runtime   string `json:"runtime"`
	Status    string `json:"status"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) GetAppInfo() AppInfo {
	info := buildinfo.Current()
	return AppInfo{
		Name:      info.Name,
		Tagline:   tagline,
		Version:   info.Version,
		Commit:    info.Commit,
		BuildTime: info.BuildTime,
		Runtime:   fmt.Sprintf("%s %s/%s", info.GoVersion, info.OS, info.Arch),
		Status:    appStatusReady,
	}
}
