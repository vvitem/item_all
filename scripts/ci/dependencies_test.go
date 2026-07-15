package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDependencyPolicyAcceptsRepositoryBaseline(t *testing.T) {
	root := dependencyFixture(t)

	if violations := checkDependencies(root); len(violations) != 0 {
		t.Fatalf("expected valid fixture, got %v", violations)
	}
}

func TestDependencyPolicyRequiresExactSQLiteVersions(t *testing.T) {
	root := dependencyFixture(t)
	if violations := checkDependencies(root); len(violations) != 0 {
		t.Fatalf("expected exact SQLite versions, got %s", violations.Error())
	}
}

func TestDependencyPolicyRejectsSQLiteVersionDrift(t *testing.T) {
	root := dependencyFixture(t)
	text, err := readText(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	text = strings.Replace(text, "modernc.org/sqlite v1.53.0", "modernc.org/sqlite v1.52.0", 1)
	writeFixture(t, filepath.Join(root, "go.mod"), text)
	assertViolation(t, checkDependencies(root), "sqlite-version")
}

func TestDependencyPolicyRejectsSQLiteVersionSuffixDrift(t *testing.T) {
	root := dependencyFixture(t)
	text, err := readText(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	text = strings.Replace(text, "modernc.org/sqlite v1.53.0", "modernc.org/sqlite v1.53.0-evil", 1)
	writeFixture(t, filepath.Join(root, "go.mod"), text)
	assertViolation(t, checkDependencies(root), "sqlite-version")
}

func TestDependencyPolicyRejectsSQLiteLibcVersionDrift(t *testing.T) {
	root := dependencyFixture(t)
	text, err := readText(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	text = strings.Replace(text, "modernc.org/libc v1.73.4", "modernc.org/libc v1.73.3", 1)
	writeFixture(t, filepath.Join(root, "go.mod"), text)
	assertViolation(t, checkDependencies(root), "sqlite-libc-version")
}

func TestDependencyPolicyRejectsFloatingNPMVersion(t *testing.T) {
	root := dependencyFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "package.json"), `{
  "packageManager": "pnpm@11.12.0",
  "engines": {"node": "24.18.0", "pnpm": "11.12.0"},
  "dependencies": {"react": "^19.1.0"},
  "devDependencies": {}
}`)

	violations := checkDependencies(root)
	assertViolation(t, violations, "npm-exact-version")
}

func TestDependencyPolicyRejectsAdditionalWorkspacePackage(t *testing.T) {
	root := dependencyFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "pnpm-workspace.yaml"), `packages:
  - .
  - ../unreviewed-package

allowBuilds:
  esbuild: true
`)

	violations := checkDependencies(root)
	assertViolation(t, violations, "pnpm-packages")
}

func TestDependencyPolicyRejectsAdditionalBuildScriptAuthorization(t *testing.T) {
	root := dependencyFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "pnpm-workspace.yaml"), `packages:
  - .

allowBuilds:
  esbuild: true
  unknown-native-package: true
`)

	violations := checkDependencies(root)
	assertViolation(t, violations, "pnpm-allow-builds")
}

func TestDependencyPolicyRejectsToolchainDrift(t *testing.T) {
	root := dependencyFixture(t)
	writeFixture(t, filepath.Join(root, "go.mod"), `module github.com/vvitem/item_all

go 1.26.0

toolchain go1.26.4

require github.com/wailsapp/wails/v2 v2.13.0
`)

	violations := checkDependencies(root)
	assertViolation(t, violations, "go-toolchain")
}

func dependencyFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFixture(t, filepath.Join(root, "go.mod"), `module github.com/vvitem/item_all

go 1.26.0

toolchain go1.26.5

require (
	github.com/wailsapp/wails/v2 v2.13.0
	modernc.org/sqlite v1.53.0
	modernc.org/libc v1.73.4 // indirect
)
`)
	writeFixture(t, filepath.Join(root, ".nvmrc"), "24.18.0\n")
	writeFixture(t, filepath.Join(root, "frontend", "package.json"), `{
  "packageManager": "pnpm@11.12.0",
  "engines": {"node": "24.18.0", "pnpm": "11.12.0"},
  "dependencies": {"react": "19.1.0"},
  "devDependencies": {"vite": "7.0.0"}
}`)
	writeFixture(t, filepath.Join(root, "frontend", "pnpm-workspace.yaml"), `packages:
  - .

allowBuilds:
  esbuild: true
`)
	return root
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertViolation(t *testing.T, violations Violations, rule string) {
	t.Helper()
	for _, violation := range violations {
		if violation.Rule == rule {
			return
		}
	}
	t.Fatalf("expected rule %q, got %s", rule, strings.TrimSpace(violations.Error()))
}
