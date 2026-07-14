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
