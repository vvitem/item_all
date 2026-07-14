package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var pinnedAction = regexp.MustCompile(`^\s*(?:-\s+)?uses:\s+([A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+)@([0-9a-f]{40})\s+#\s+v[0-9][^\s]*\s*$`)
var writePermission = regexp.MustCompile(`(?m)^\s+[A-Za-z][A-Za-z-]*:\s*write\s*$`)

var requiredActions = map[string]string{
	"actions/checkout":   "9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0",
	"actions/setup-go":   "924ae3a1cded613372ab5595356fb5720e22ba16",
	"actions/setup-node": "820762786026740c76f36085b0efc47a31fe5020",
	"actions/cache":      "55cc8345863c7cc4c66a329aec7e433d2d1c52a9",
}

func checkWorkflow(root string) Violations {
	var violations Violations
	path := filepath.Join(root, ".github", "workflows", "ci.yml")
	text, err := readText(path)
	if err != nil {
		violations.Add(".github/workflows/ci.yml", "workflow-readable", err.Error())
		return violations
	}
	mustContain := map[string]string{
		"trigger-pr":         "  pull_request:",
		"trigger-push":       "  push:",
		"trigger-manual":     "  workflow_dispatch:",
		"permissions":        "  contents: read",
		"concurrency-group":  "group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}",
		"concurrency-cancel": "cancel-in-progress: true",
	}
	for rule, token := range mustContain {
		if !strings.Contains(text, token) {
			violations.Add(".github/workflows/ci.yml", rule, fmt.Sprintf("expected %q", token))
		}
	}
	forbidden := map[string]string{
		"pull-request-target": "pull_request_target",
		"path-filter":         "paths:",
		"path-ignore":         "paths-ignore:",
		"continue-on-error":   "continue-on-error:",
		"ignored-exit-code":   "|| true",
		"write-all":           "write-all",
	}
	for rule, token := range forbidden {
		if strings.Contains(text, token) {
			violations.Add(".github/workflows/ci.yml", rule, fmt.Sprintf("forbidden token %q", token))
		}
	}
	if writePermission.MatchString(text) {
		violations.Add(".github/workflows/ci.yml", "permissions", "workflow and jobs must not request write permissions")
	}

	expectedJobs := map[string]struct {
		runner  string
		timeout string
		depth   string
	}{
		"quality":                {runner: "ubuntu-24.04", timeout: "15", depth: "1"},
		"generated-and-security": {runner: "ubuntu-24.04", timeout: "15", depth: "0"},
		"windows-build":          {runner: "windows-2025", timeout: "25", depth: "1"},
	}
	for job, expectation := range expectedJobs {
		block := workflowJobBlock(text, job)
		if block == "" {
			violations.Add(".github/workflows/ci.yml", "job-"+job, "required job is missing")
			continue
		}
		required := []string{
			"name: " + job,
			"runs-on: " + expectation.runner,
			"timeout-minutes: " + expectation.timeout,
			"fetch-depth: " + expectation.depth,
			"persist-credentials: false",
		}
		for _, token := range required {
			if !strings.Contains(block, token) {
				rule := "job-" + job
				if job == "generated-and-security" && strings.HasPrefix(token, "fetch-depth:") {
					rule = "security-history"
				}
				violations.Add(".github/workflows/ci.yml", rule, fmt.Sprintf("expected %q", token))
			}
		}
	}

	for _, line := range strings.Split(text, "\n") {
		if !strings.Contains(line, "uses:") {
			continue
		}
		match := pinnedAction.FindStringSubmatch(line)
		if match == nil {
			violations.Add(".github/workflows/ci.yml", "action-sha", "every uses entry must use a 40-character SHA and version comment")
			continue
		}
		if expected, ok := requiredActions[match[1]]; ok && match[2] != expected {
			violations.Add(".github/workflows/ci.yml", "action-version", fmt.Sprintf("%s must use %s", match[1], expected))
		}
	}
	checkUntrustedRunContext(text, &violations)
	return violations
}

func workflowJobBlock(text, job string) string {
	lines := strings.Split(text, "\n")
	marker := "  " + job + ":"
	start := -1
	for index, line := range lines {
		if line == marker {
			start = index
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := len(lines)
	for index := start + 1; index < len(lines); index++ {
		line := lines[index]
		if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") && strings.HasSuffix(line, ":") {
			end = index
			break
		}
	}
	return strings.Join(lines[start:end], "\n")
}

func checkUntrustedRunContext(text string, violations *Violations) {
	risky := []string{
		"github.event.pull_request.title",
		"github.event.pull_request.body",
		"github.event.pull_request.head.label",
		"github.head_ref",
		"github.event.issue",
	}
	lines := strings.Split(text, "\n")
	inRun := false
	runIndent := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if strings.HasPrefix(trimmed, "run:") {
			inRun = true
			runIndent = indent
		} else if inRun && trimmed != "" && indent <= runIndent {
			inRun = false
		}
		if !inRun {
			continue
		}
		for _, token := range risky {
			if strings.Contains(line, token) {
				violations.Add(".github/workflows/ci.yml", "untrusted-context", fmt.Sprintf("run block contains %s", token))
			}
		}
	}
}
