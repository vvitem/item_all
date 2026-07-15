package main

import (
	"fmt"

	"github.com/vvitem/item_all/internal/buildinfo"
	"github.com/vvitem/item_all/internal/storage"
)

const tagline = "Local-first SafeOps"
const appStatusReady = "ready"
const appStatusStorageError = "storage_error"

type App struct {
	startupErr error
}

type AppInfo struct {
	Name      string             `json:"name"`
	Tagline   string             `json:"tagline"`
	Version   string             `json:"version"`
	Commit    string             `json:"commit"`
	BuildTime string             `json:"buildTime"`
	Runtime   string             `json:"runtime"`
	Status    string             `json:"status"`
	Error     *storage.SafeError `json:"error,omitempty"`
}

func NewApp(startupErr error) *App {
	return &App{startupErr: startupErr}
}

func (a *App) GetAppInfo() AppInfo {
	info := buildinfo.Current()
	result := AppInfo{
		Name:      info.Name,
		Tagline:   tagline,
		Version:   info.Version,
		Commit:    info.Commit,
		BuildTime: info.BuildTime,
		Runtime:   fmt.Sprintf("%s %s/%s", info.GoVersion, info.OS, info.Arch),
		Status:    appStatusReady,
	}
	if a != nil && a.startupErr != nil {
		safe := storage.Project(a.startupErr)
		result.Status = appStatusStorageError
		result.Error = &safe
	}
	return result
}
