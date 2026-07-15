package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	expectedModule            = "github.com/vvitem/item_all"
	expectedGoDirective       = "1.26.0"
	expectedGoToolchain       = "go1.26.5"
	expectedWailsVersion      = "v2.13.0"
	expectedNodeVersion       = "24.18.0"
	expectedPNPMVersion       = "11.12.0"
	expectedSQLiteVersion     = "v1.53.0"
	expectedSQLiteLibcVersion = "v1.73.4"
)

var exactVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)

type packageManifest struct {
	PackageManager  string            `json:"packageManager"`
	Engines         map[string]string `json:"engines"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func checkDependencies(root string) Violations {
	var violations Violations
	checkGoModule(root, &violations)
	checkNodeManifest(root, &violations)
	checkPNPMWorkspace(root, &violations)
	return violations
}

func checkGoModule(root string, violations *Violations) {
	path := filepath.Join(root, "go.mod")
	text, err := readText(path)
	if err != nil {
		violations.Add("go.mod", "go-mod-readable", err.Error())
		return
	}
	required := map[string]string{
		"go-module":     "module " + expectedModule,
		"go-version":    "go " + expectedGoDirective,
		"go-toolchain":  "toolchain " + expectedGoToolchain,
		"wails-version": "github.com/wailsapp/wails/v2 " + expectedWailsVersion,
	}
	for rule, token := range required {
		if !strings.Contains(text, token) {
			violations.Add("go.mod", rule, fmt.Sprintf("expected %q", token))
		}
	}
	checkExactModuleVersion(text, "modernc.org/sqlite", expectedSQLiteVersion, "sqlite-version", violations)
	checkExactModuleVersion(text, "modernc.org/libc", expectedSQLiteLibcVersion, "sqlite-libc-version", violations)
}

func checkExactModuleVersion(text, module, version, rule string, violations *Violations) {
	pattern := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(module) + `\s+` + regexp.QuoteMeta(version) + `(?:\s+//[^\n]*)?\s*$`)
	if !pattern.MatchString(text) {
		violations.Add("go.mod", rule, fmt.Sprintf("expected exact module version %s %s", module, version))
	}
}

func checkNodeManifest(root string, violations *Violations) {
	path := filepath.Join(root, "frontend", "package.json")
	text, err := readText(path)
	if err != nil {
		violations.Add("frontend/package.json", "package-readable", err.Error())
		return
	}
	var manifest packageManifest
	if err := json.Unmarshal([]byte(text), &manifest); err != nil {
		violations.Add("frontend/package.json", "package-json", err.Error())
		return
	}
	if manifest.PackageManager != "pnpm@"+expectedPNPMVersion {
		violations.Add("frontend/package.json", "pnpm-version", "packageManager must be pnpm@"+expectedPNPMVersion)
	}
	if manifest.Engines["node"] != expectedNodeVersion {
		violations.Add("frontend/package.json", "node-version", "engines.node must match .nvmrc")
	}
	if manifest.Engines["pnpm"] != expectedPNPMVersion {
		violations.Add("frontend/package.json", "pnpm-engine", "engines.pnpm must be exact")
	}
	nvmrc, err := readText(filepath.Join(root, ".nvmrc"))
	if err != nil || strings.TrimSpace(nvmrc) != expectedNodeVersion {
		violations.Add(".nvmrc", "node-version", "expected "+expectedNodeVersion)
	}
	names := make([]string, 0, len(manifest.Dependencies)+len(manifest.DevDependencies))
	for name := range manifest.Dependencies {
		names = append(names, name)
	}
	for name := range manifest.DevDependencies {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		version := manifest.Dependencies[name]
		if version == "" {
			version = manifest.DevDependencies[name]
		}
		if !exactVersion.MatchString(version) {
			violations.Add("frontend/package.json", "npm-exact-version", fmt.Sprintf("%s uses non-exact version %q", name, version))
		}
	}
}

func checkPNPMWorkspace(root string, violations *Violations) {
	path := filepath.Join(root, "frontend", "pnpm-workspace.yaml")
	text, err := readText(path)
	if err != nil {
		violations.Add("frontend/pnpm-workspace.yaml", "pnpm-workspace-readable", err.Error())
		return
	}
	lines := strings.Split(text, "\n")
	packages := parseYAMLListSection(lines, "packages:")
	if len(packages) != 1 || packages[0] != "." {
		violations.Add("frontend/pnpm-workspace.yaml", "pnpm-packages", "packages must contain only the frontend package root: .")
	}

	inAllowBuilds := false
	allowed := map[string]bool{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "allowBuilds:" {
			inAllowBuilds = true
			continue
		}
		if !inAllowBuilds {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(line, "  ") {
			inAllowBuilds = false
			continue
		}
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[1]) == "true" {
			allowed[parts[0]] = true
		}
	}
	if len(allowed) != 1 || !allowed["esbuild"] {
		violations.Add("frontend/pnpm-workspace.yaml", "pnpm-allow-builds", "allowBuilds must contain only esbuild: true")
	}
}

func parseYAMLListSection(lines []string, header string) []string {
	inSection := false
	values := []string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == header {
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		if trimmed == "" {
			continue
		}
		if !strings.HasPrefix(line, "  - ") {
			break
		}
		value := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-"))
		values = append(values, value)
	}
	return values
}
