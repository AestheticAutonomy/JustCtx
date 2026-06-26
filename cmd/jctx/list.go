package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List metadata and dynamic dimensions",
}

type ProviderJSON struct {
	Name           string   `json:"name"`
	SupportedTypes []string `json:"supported_types"`
}

var listProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List all registered providers and their supported types",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runListProviders(os.Stdout, jsonFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runListProviders(w io.Writer, jsonOutput bool) error {
	allProviders := providers.All()
	sort.Slice(allProviders, func(i, j int) bool {
		return allProviders[i].Name() < allProviders[j].Name()
	})

	if jsonOutput {
		var out []ProviderJSON
		for _, p := range allProviders {
			var supported []string
			for _, t := range p.SupportedTypes() {
				supported = append(supported, string(t))
			}
			out = append(out, ProviderJSON{
				Name:           p.Name(),
				SupportedTypes: supported,
			})
		}
		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(data))
		return nil
	}

	for _, p := range allProviders {
		var supported []string
		for _, t := range p.SupportedTypes() {
			supported = append(supported, string(t))
		}
		fmt.Fprintf(w, "%-12s %s\n", p.Name(), strings.Join(supported, " "))
	}
	return nil
}

type TypeJSON struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

var allTypes = []TypeJSON{
	{string(schema.TypeRules), "AI coding guidelines and instructions"},
	{string(schema.TypeIgnore), "Files and patterns to exclude"},
	{string(schema.TypeMCP), "MCP server configuration"},
	{string(schema.TypeCommands), "Slash commands / custom prompts"},
	{string(schema.TypeHooks), "Event hooks"},
	{string(schema.TypeSkills), "Reusable skill definitions"},
	{string(schema.TypeSubagents), "Subagent definitions"},
	{string(schema.TypePermissions), "Permission rules"},
	{string(schema.TypePolicies), "Policy rules (reserved)"},
}

var listTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List all configuration types and their descriptions",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runListTypes(os.Stdout, jsonFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runListTypes(w io.Writer, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(allTypes, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(data))
		return nil
	}

	for _, t := range allTypes {
		fmt.Fprintf(w, "%-12s %s\n", t.Type, t.Description)
	}
	return nil
}

type TagJSON struct {
	Tag string `json:"tag"`
}

var listTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all unique tags used in the project source files",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}
		if err := runListTags(cwd, os.Stdout, jsonFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runListTags(cwd string, w io.Writer, jsonOutput bool) error {
	jctxDir, err := findJctxDir(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "no .jctx/ directory found")
		return nil
	}

	tagsMap := make(map[string]bool)
	err = filepath.Walk(jctxDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() != ".jctx" && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		vals, err := parseFrontmatterValues(path, "tag")
		if err == nil {
			for _, val := range vals {
				tagsMap[val] = true
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	var tags []string
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	if jsonOutput {
		var out []TagJSON
		for _, tag := range tags {
			out = append(out, TagJSON{Tag: tag})
		}
		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(data))
		return nil
	}

	for _, tag := range tags {
		fmt.Fprintln(w, tag)
	}
	return nil
}

type RoleJSON struct {
	Role string `json:"role"`
}

var listRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List all unique roles used in the project source files",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}
		if err := runListRoles(cwd, os.Stdout, jsonFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runListRoles(cwd string, w io.Writer, jsonOutput bool) error {
	jctxDir, err := findJctxDir(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "no .jctx/ directory found")
		return nil
	}

	rolesMap := make(map[string]bool)
	err = filepath.Walk(jctxDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() != ".jctx" && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		vals, err := parseFrontmatterValues(path, "role")
		if err == nil {
			for _, val := range vals {
				rolesMap[val] = true
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	var roles []string
	for role := range rolesMap {
		roles = append(roles, role)
	}
	sort.Strings(roles)

	if jsonOutput {
		var out []RoleJSON
		for _, role := range roles {
			out = append(out, RoleJSON{Role: role})
		}
		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(data))
		return nil
	}

	for _, role := range roles {
		fmt.Fprintln(w, role)
	}
	return nil
}

func init() {
	listCmd.AddCommand(listProvidersCmd)
	listCmd.AddCommand(listTypesCmd)
	listCmd.AddCommand(listTagsCmd)
	listCmd.AddCommand(listRolesCmd)
	rootCmd.AddCommand(listCmd)
}

func findRepoRoot(start string) string {
	dir, err := filepath.Abs(start)
	if err != nil {
		return start
	}
	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return start
}

func findJctxDir(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}
	for {
		jctxPath := filepath.Join(dir, ".jctx")
		if info, err := os.Stat(jctxPath); err == nil && info.IsDir() {
			return jctxPath, nil
		}

		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no .jctx/ directory found")
}

func extractFrontmatter(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return ""
	}
	lines := strings.Split(content, "\n")
	firstIdx := -1
	secondIdx := -1
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "---" {
			if firstIdx == -1 {
				firstIdx = i
			} else {
				secondIdx = i
				break
			}
		}
	}
	if firstIdx != -1 && secondIdx != -1 && secondIdx > firstIdx {
		return strings.Join(lines[firstIdx+1:secondIdx], "\n")
	}
	return ""
}

func parseFrontmatterValues(filePath string, valueKey string) ([]string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(bytes)
	ext := strings.ToLower(filepath.Ext(filePath))
	var yamlStr string

	if ext == ".md" || ext == ".mdc" {
		yamlStr = extractFrontmatter(content)
		if yamlStr == "" {
			return nil, nil
		}
	} else if ext == ".yaml" || ext == ".yml" || ext == ".json" {
		yamlStr = content
	} else {
		return nil, nil
	}

	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
		return nil, err
	}

	var found []string
	keysToSearch := []string{valueKey}
	if valueKey == "tag" {
		keysToSearch = append(keysToSearch, "tags")
	} else if valueKey == "role" {
		keysToSearch = append(keysToSearch, "roles")
	}

	for _, key := range keysToSearch {
		findValues(data, key, &found)
	}

	return found, nil
}

func findValues(val interface{}, keyName string, collect *[]string) {
	if val == nil {
		return
	}
	switch v := val.(type) {
	case map[string]interface{}:
		for k, child := range v {
			if strings.EqualFold(k, keyName) {
				collectValues(child, collect)
			} else {
				findValues(child, keyName, collect)
			}
		}
	case []interface{}:
		for _, child := range v {
			findValues(child, keyName, collect)
		}
	}
}

func collectValues(val interface{}, collect *[]string) {
	if val == nil {
		return
	}
	switch v := val.(type) {
	case string:
		if v != "" {
			*collect = append(*collect, v)
		}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				*collect = append(*collect, s)
			}
		}
	case map[string]interface{}:
		for _, child := range v {
			collectValues(child, collect)
		}
	}
}
