package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunListProviders(t *testing.T) {
	var buf bytes.Buffer
	err := runListProviders(&buf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	expectedLines := []string{
		"agents       rules",
		"antigravity  rules",
		"claude       rules",
		"cursor       rules",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, but got:\n%s", expected, output)
		}
	}
}

func TestRunListProviders_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := runListProviders(&buf, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var providers []ProviderJSON
	if err := json.Unmarshal(buf.Bytes(), &providers); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v, raw output: %s", err, buf.String())
	}

	if len(providers) == 0 {
		t.Fatal("expected registered providers, got 0")
	}

	foundClaude := false
	for _, p := range providers {
		if p.Name == "claude" {
			foundClaude = true
			if len(p.SupportedTypes) != 1 || p.SupportedTypes[0] != "rules" {
				t.Errorf("claude supported types mismatch: %v", p.SupportedTypes)
			}
		}
	}

	if !foundClaude {
		t.Error("claude provider not found in JSON output")
	}
}

func TestRunListTypes(t *testing.T) {
	var buf bytes.Buffer
	err := runListTypes(&buf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	expectedSubstring := "rules        AI coding guidelines and instructions"
	if !strings.Contains(output, expectedSubstring) {
		t.Errorf("expected output to contain %q, but got:\n%s", expectedSubstring, output)
	}
}

func TestRunListTypes_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := runListTypes(&buf, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var types []TypeJSON
	if err := json.Unmarshal(buf.Bytes(), &types); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v, raw output: %s", err, buf.String())
	}

	foundRules := false
	for _, typ := range types {
		if typ.Type == "rules" {
			foundRules = true
			expectedDesc := "AI coding guidelines and instructions"
			if typ.Description != expectedDesc {
				t.Errorf("rules description mismatch: expected %q, got %q", expectedDesc, typ.Description)
			}
		}
	}

	if !foundRules {
		t.Error("rules type not found in JSON output")
	}
}

func TestRunListTagsAndRoles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a dummy .git file to bound the walk
	gitPath := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitPath, 0755); err != nil {
		t.Fatalf("failed to create temporary .git: %v", err)
	}

	// Create .jctx subfolders
	jctxDir := filepath.Join(tmpDir, ".jctx")
	if err := os.MkdirAll(filepath.Join(jctxDir, "rules"), 0755); err != nil {
		t.Fatalf("failed to create rules folder: %v", err)
	}

	// Write a rules markdown file with YAML frontmatter
	ruleContent := `---
placement: project_root
targets: [claude, cursor]
role: dbreviewer
tags: [postgres, go]
---
# Coding guidelines
`
	if err := os.WriteFile(filepath.Join(jctxDir, "rules", "coding-guidelines.md"), []byte(ruleContent), 0644); err != nil {
		t.Fatalf("failed to write coding-guidelines: %v", err)
	}

	// Write a hook YAML file
	hookContent := `targets: [claude]
hooks:
  - event: PreToolUse
    command: ./check.sh
    role: security
    tags: [audit]
`
	if err := os.WriteFile(filepath.Join(jctxDir, "hooks.yaml"), []byte(hookContent), 0644); err != nil {
		t.Fatalf("failed to write hooks.yaml: %v", err)
	}

	// Test runListTags
	var tagsBuf bytes.Buffer
	err := runListTags(tmpDir, &tagsBuf, false)
	if err != nil {
		t.Fatalf("runListTags unexpected error: %v", err)
	}

	tagsOutput := strings.TrimSpace(tagsBuf.String())
	expectedTags := []string{"audit", "go", "postgres"}
	actualTags := strings.Split(strings.ReplaceAll(tagsOutput, "\r\n", "\n"), "\n")
	if len(actualTags) != len(expectedTags) {
		t.Fatalf("expected tags count %d, got %d. Output: %q", len(expectedTags), len(actualTags), actualTags)
	}
	for i, tag := range expectedTags {
		if actualTags[i] != tag {
			t.Errorf("at index %d: expected tag %q, got %q", i, tag, actualTags[i])
		}
	}

	// Test runListTags JSON
	var tagsJSONBuf bytes.Buffer
	err = runListTags(tmpDir, &tagsJSONBuf, true)
	if err != nil {
		t.Fatalf("runListTags JSON unexpected error: %v", err)
	}

	var tagsJSON []TagJSON
	if err := json.Unmarshal(tagsJSONBuf.Bytes(), &tagsJSON); err != nil {
		t.Fatalf("failed to unmarshal JSON tags: %v", err)
	}

	if len(tagsJSON) != len(expectedTags) {
		t.Fatalf("expected %d tags in JSON, got %d", len(expectedTags), len(tagsJSON))
	}
	for i, tag := range expectedTags {
		if tagsJSON[i].Tag != tag {
			t.Errorf("at index %d: expected tag %q, got %q", i, tag, tagsJSON[i].Tag)
		}
	}

	// Test runListRoles
	var rolesBuf bytes.Buffer
	err = runListRoles(tmpDir, &rolesBuf, false)
	if err != nil {
		t.Fatalf("runListRoles unexpected error: %v", err)
	}

	rolesOutput := strings.TrimSpace(rolesBuf.String())
	expectedRoles := []string{"dbreviewer", "security"}
	actualRoles := strings.Split(strings.ReplaceAll(rolesOutput, "\r\n", "\n"), "\n")
	if len(actualRoles) != len(expectedRoles) {
		t.Fatalf("expected roles count %d, got %d. Output: %q", len(expectedRoles), len(actualRoles), actualRoles)
	}
	for i, role := range expectedRoles {
		if actualRoles[i] != role {
			t.Errorf("at index %d: expected role %q, got %q", i, role, actualRoles[i])
		}
	}
}
