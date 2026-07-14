package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Violation struct {
	Path    string
	Rule    string
	Message string
}

type Violations []Violation

func (v *Violations) Add(path, rule, message string) {
	*v = append(*v, Violation{Path: filepath.ToSlash(path), Rule: rule, Message: message})
}

func (v Violations) Error() string {
	if len(v) == 0 {
		return ""
	}
	sort.Slice(v, func(i, j int) bool {
		if v[i].Path != v[j].Path {
			return v[i].Path < v[j].Path
		}
		return v[i].Rule < v[j].Rule
	})
	var builder strings.Builder
	for _, item := range v {
		fmt.Fprintf(&builder, "%s [%s] %s\n", item.Path, item.Rule, item.Message)
	}
	return builder.String()
}

func runPolicy(name, root string) error {
	var violations Violations
	switch name {
	case "dependencies":
		violations = append(violations, checkDependencies(root)...)
	case "workflow":
		violations = append(violations, checkWorkflow(root)...)
	case "boundary":
		violations = append(violations, checkBoundary(root)...)
	case "all":
		violations = append(violations, checkDependencies(root)...)
		violations = append(violations, checkWorkflow(root)...)
		violations = append(violations, checkBoundary(root)...)
	default:
		return fmt.Errorf("unknown policy %q; use dependencies, workflow, boundary, or all", name)
	}
	if len(violations) != 0 {
		return violations
	}
	return nil
}

func findRepositoryRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if fileExists(filepath.Join(current, "go.mod")) && fileExists(filepath.Join(current, "frontend", "package.json")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", errors.New("repository root not found")
		}
		current = parent
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func readText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
