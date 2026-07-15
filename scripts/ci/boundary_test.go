package main

import (
	"path/filepath"
	"testing"
)

func TestBoundaryPolicyAcceptsMinimalApplication(t *testing.T) {
	root := boundaryFixture(t)

	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected valid boundary, got %s", violations.Error())
	}
}

func TestBoundaryPolicyAcceptsJavaScriptLineComment(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "src", "comment.ts"), "// local implementation note\nexport const ready = true\n")

	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected line comment to be valid, got %s", violations.Error())
	}
}

func TestBoundaryPolicyRejectsPrivateKeyFile(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "test.key"), "not-a-real-key")

	assertViolation(t, checkBoundary(root), "secret-file")
}

func TestBoundaryPolicyRejectsRemoteFrontendResource(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "src", "styles.css"), `@import url("https://example.invalid/font.css");`)

	assertViolation(t, checkBoundary(root), "remote-resource")
}

func TestBoundaryPolicyRejectsProtocolRelativeFrontendResource(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "src", "loader.ts"), `const source = "//cdn.example.invalid/app.js"`)

	assertViolation(t, checkBoundary(root), "remote-resource")
}

func TestBoundaryPolicyRejectsEnvironmentRead(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "app.go"), "package main\nimport \"os\"\nvar value = os.Getenv(\"TOKEN\")\n")

	assertViolation(t, checkBoundary(root), "environment-read")
}

func TestBoundaryPolicyAllowsDatabaseSQLOnlyInsideStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "querier.go"), "package storage\nimport \"database/sql\"\nvar _ sql.Result\n")
	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected storage import to pass, got %s", violations.Error())
	}
}

func TestBoundaryPolicyRejectsDatabaseSQLOutsideStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "app.go"), "package main\nimport \"database/sql\"\nvar _ sql.Result\n")
	assertViolation(t, checkBoundary(root), "scope-import")
}

func TestBoundaryPolicyAllowsModernSQLiteOnlyInsideSQLiteStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "sqlite", "store.go"), "package sqlite\nimport _ \"modernc.org/sqlite\"\n")
	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected driver import to pass, got %s", violations.Error())
	}
}

func TestBoundaryPolicyRejectsModernSQLiteOutsideSQLiteStorage(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "store.go"), "package storage\nimport _ \"modernc.org/sqlite\"\n")
	assertViolation(t, checkBoundary(root), "scope-import")
}

func TestBoundaryPolicyAllowsOnlyStandardXDGDataHomeLookup(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "sqlite", "path_linux.go"), "package sqlite\nimport \"os\"\nfunc root() string { v, _ := os.LookupEnv(\"XDG_DATA_HOME\"); return v }\n")
	if violations := checkBoundary(root); len(violations) != 0 {
		t.Fatalf("expected XDG lookup to pass, got %s", violations.Error())
	}
}

func TestBoundaryPolicyRejectsOtherEnvironmentReadInLinuxPath(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "internal", "storage", "sqlite", "path_linux.go"), "package sqlite\nimport \"os\"\nfunc root() string { return os.Getenv(\"HOME\") }\n")
	assertViolation(t, checkBoundary(root), "environment-read")
}

func TestBoundaryPolicyRejectsAdditionalWailsBinding(t *testing.T) {
	root := boundaryFixture(t)
	writeFixture(t, filepath.Join(root, "frontend", "wailsjs", "go", "main", "App.js"), `export function GetAppInfo() {}
export function ReadFile() {}
`)

	assertViolation(t, checkBoundary(root), "wails-binding")
}

func boundaryFixture(t *testing.T) string {
	t.Helper()
	root := dependencyFixture(t)
	writeFixture(t, filepath.Join(root, "app.go"), "package main\n")
	writeFixture(t, filepath.Join(root, "frontend", "src", "App.tsx"), "export default function App() { return null }\n")
	writeFixture(t, filepath.Join(root, "frontend", "wailsjs", "go", "main", "App.js"), "export function GetAppInfo() {}\n")
	return root
}
