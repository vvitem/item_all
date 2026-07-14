package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowPolicyAcceptsLayeredBaseline(t *testing.T) {
	root := dependencyFixture(t)
	writeFixture(t, filepath.Join(root, ".github", "workflows", "ci.yml"), validWorkflowFixture)

	if violations := checkWorkflow(root); len(violations) != 0 {
		t.Fatalf("expected valid workflow, got %s", violations.Error())
	}
}

func TestWorkflowPolicyRejectsMutableActionTag(t *testing.T) {
	root := dependencyFixture(t)
	invalid := strings.Replace(validWorkflowFixture,
		"actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0",
		"actions/checkout@v7",
		1,
	)
	writeFixture(t, filepath.Join(root, ".github", "workflows", "ci.yml"), invalid)

	assertViolation(t, checkWorkflow(root), "action-sha")
}

func TestWorkflowPolicyRejectsWritePermission(t *testing.T) {
	root := dependencyFixture(t)
	invalid := strings.Replace(validWorkflowFixture, "contents: read", "contents: write", 1)
	writeFixture(t, filepath.Join(root, ".github", "workflows", "ci.yml"), invalid)

	assertViolation(t, checkWorkflow(root), "permissions")
}

func TestWorkflowPolicyRejectsShallowSecurityCheckout(t *testing.T) {
	root := dependencyFixture(t)
	invalid := strings.Replace(validWorkflowFixture, "fetch-depth: 0", "fetch-depth: 1", 1)
	writeFixture(t, filepath.Join(root, ".github", "workflows", "ci.yml"), invalid)

	assertViolation(t, checkWorkflow(root), "security-history")
}

func TestWorkflowPolicyRejectsPullRequestTarget(t *testing.T) {
	root := dependencyFixture(t)
	invalid := strings.Replace(validWorkflowFixture, "pull_request:", "pull_request_target:", 1)
	writeFixture(t, filepath.Join(root, ".github", "workflows", "ci.yml"), invalid)

	assertViolation(t, checkWorkflow(root), "pull-request-target")
}

const validWorkflowFixture = `name: CI
on:
  pull_request:
    branches: [main]
  push:
    branches: [main]
  workflow_dispatch:
permissions:
  contents: read
concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true
jobs:
  quality:
    name: quality
    runs-on: ubuntu-24.04
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          fetch-depth: 1
          persist-credentials: false
      - uses: actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16 # v6.5.0
      - uses: actions/setup-node@820762786026740c76f36085b0efc47a31fe5020 # v7.0.0
      - uses: actions/cache@55cc8345863c7cc4c66a329aec7e433d2d1c52a9 # v6.1.0
  generated-and-security:
    name: generated-and-security
    runs-on: ubuntu-24.04
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          fetch-depth: 0
          persist-credentials: false
      - uses: actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16 # v6.5.0
      - uses: actions/setup-node@820762786026740c76f36085b0efc47a31fe5020 # v7.0.0
      - uses: actions/cache@55cc8345863c7cc4c66a329aec7e433d2d1c52a9 # v6.1.0
  windows-build:
    name: windows-build
    runs-on: windows-2025
    timeout-minutes: 25
    steps:
      - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7.0.0
        with:
          fetch-depth: 1
          persist-credentials: false
      - uses: actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16 # v6.5.0
      - uses: actions/setup-node@820762786026740c76f36085b0efc47a31fe5020 # v7.0.0
      - uses: actions/cache@55cc8345863c7cc4c66a329aec7e433d2d1c52a9 # v6.1.0
`
