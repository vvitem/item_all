package main

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var absoluteRemoteResource = regexp.MustCompile(`(?i)https?://[^\s"')>]+`)
var protocolRelativeResource = regexp.MustCompile(`(?i)["'(=][ \t]*//[a-z0-9][a-z0-9.-]*(:[0-9]+)?(/[^\s"')>]*)?`)
var bindingExport = regexp.MustCompile(`(?m)^export function ([A-Za-z0-9_]+)\(`)

func checkBoundary(root string) Violations {
	var violations Violations
	forbiddenPaths := []string{
		"internal/ssh",
		"internal/sftp",
		"internal/database",
		"internal/operation",
		"internal/ai",
		"internal/credential",
		"internal/telemetry",
	}
	forbiddenNames := map[string]bool{
		".env":       true,
		".env.local": true,
		"id_rsa":     true,
		"id_ed25519": true,
	}
	forbiddenExtensions := map[string]bool{
		".pem": true,
		".key": true,
		".p12": true,
		".pfx": true,
	}

	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			violations.Add(relative(root, path), "walk", walkErr.Error())
			return nil
		}
		rel := relative(root, path)
		if entry.IsDir() {
			if rel == ".git" || rel == "frontend/node_modules" || rel == "frontend/dist" || rel == "build/bin" {
				return filepath.SkipDir
			}
			for _, forbidden := range forbiddenPaths {
				if rel == forbidden || strings.HasPrefix(rel, forbidden+"/") {
					violations.Add(rel, "scope-directory", "capability is outside the currently approved task scope")
				}
			}
			return nil
		}
		lowerName := strings.ToLower(entry.Name())
		if forbiddenNames[lowerName] || forbiddenExtensions[strings.ToLower(filepath.Ext(lowerName))] {
			violations.Add(rel, "secret-file", "secret-bearing file type is forbidden")
		}
		if shouldScanFrontend(rel) || shouldScanSource(rel) {
			text, err := readText(path)
			if err != nil {
				violations.Add(rel, "source-readable", err.Error())
				return nil
			}
			if shouldScanFrontend(rel) && containsRemoteResource(text) {
				violations.Add(rel, "remote-resource", "frontend source must not load remote scripts, styles, fonts, or CDN resources")
			}
			if shouldScanSource(rel) {
				checkEnvironmentReads(rel, text, &violations)
				checkCapabilityImports(rel, text, &violations)
			}
		}
		return nil
	})
	checkBindings(root, &violations)
	return violations
}

func checkEnvironmentReads(rel, text string, violations *Violations) {
	sanitized := text
	if rel == "internal/storage/sqlite/path_linux.go" {
		sanitized = strings.ReplaceAll(sanitized, `os.LookupEnv("XDG_DATA_HOME")`, "")
	}
	for _, token := range []string{"os.Getenv(", "os.LookupEnv(", "process.env", "import.meta.env"} {
		if strings.Contains(sanitized, token) {
			violations.Add(rel, "environment-read", "only XDG_DATA_HOME in the Linux storage path resolver is approved")
		}
	}
}

func checkCapabilityImports(rel, text string, violations *Violations) {
	rules := []struct {
		token         string
		allowedPrefix string
	}{
		{token: "golang.org/x/crypto/ssh"},
		{token: "database/sql", allowedPrefix: "internal/storage/"},
		{token: "modernc.org/sqlite", allowedPrefix: "internal/storage/sqlite/"},
		{token: "github.com/openai/"},
		{token: "go.opentelemetry.io/"},
	}
	for _, rule := range rules {
		if !strings.Contains(text, rule.token) {
			continue
		}
		if rule.allowedPrefix != "" && strings.HasPrefix(rel, rule.allowedPrefix) {
			continue
		}
		violations.Add(rel, "scope-import", "capability dependency is outside its approved package boundary: "+rule.token)
	}
}

func containsRemoteResource(text string) bool {
	return absoluteRemoteResource.MatchString(text) || protocolRelativeResource.MatchString(text)
}

func checkBindings(root string, violations *Violations) {
	path := filepath.Join(root, "frontend", "wailsjs", "go", "main", "App.js")
	text, err := readText(path)
	if err != nil {
		violations.Add("frontend/wailsjs/go/main/App.js", "wails-binding", err.Error())
		return
	}
	matches := bindingExport.FindAllStringSubmatch(text, -1)
	names := make([]string, 0, len(matches))
	for _, match := range matches {
		names = append(names, match[1])
	}
	sort.Strings(names)
	if len(names) != 1 || names[0] != "GetAppInfo" {
		violations.Add("frontend/wailsjs/go/main/App.js", "wails-binding", "only GetAppInfo may be exported")
	}
}

func relative(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

func shouldScanFrontend(rel string) bool {
	return rel == "frontend/index.html" || strings.HasPrefix(rel, "frontend/src/")
}

func shouldScanSource(rel string) bool {
	if strings.HasPrefix(rel, "scripts/ci/") || strings.HasPrefix(rel, "frontend/wailsjs/") {
		return false
	}
	extension := strings.ToLower(filepath.Ext(rel))
	return extension == ".go" || extension == ".ts" || extension == ".tsx" || extension == ".js" || extension == ".jsx"
}
