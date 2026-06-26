# Cursor Provider Specification

## File Conventions

### 1. Primary Format: `.cursor/rules/*.mdc`
- **Location:** `.cursor/rules/` under the project root.
- **Structure:** MDC rules are markdown files ending in `.mdc` with a YAML frontmatter block containing metadata at the very top:
  ```yaml
  ---
  description: "Rule description used by the LLM for request-based activation"
  globs: ["glob patterns (e.g. *.go, internal/**/*.go)"]
  alwaysApply: false
  ---
  ```
- **Frontmatter Fields:**
  - `description` (String): **Critical.** Acts as a semantic prompt for Cursor's routing system to dynamically auto-attach the rule based on the context of the user request.
  - `globs` (Array of Strings or String): Specifies which files or directories match this rule to trigger auto-attachment when those files are active or being edited.
  - `alwaysApply` (Boolean): If `true`, the rule is persistent across every single conversation (token-heavy, use sparingly).
- **Plain Markdown:** Plain `.md` files located in `.cursor/rules/` are passive documentation files. They are NOT treated as active rules and are ignored for automatic rule enforcement (frontmatter and `.mdc` extension are required).

### 2. Model Context Protocol (MCP)
- **Project-Scoped:** `.cursor/mcp.json` at the project root.
- **User-Scoped:** `~/.cursor/mcp.json` (Windows: `%USERPROFILE%\.cursor\mcp.json`), with servers defined under the `mcpServers` field.

### 3. Saved Plans (Agent Mode)
- **Location:** `.cursor/plans/` directory under the project root. Used to store and track agent plans and roadmaps if they are saved to the workspace.

### 4. Activation Modes
Cursor uses four activation modes for MDC rules:
- **Always Apply:** `alwaysApply: true` in frontmatter.
- **Auto Attached:** Driven by the `globs` match when working on relevant files.
- **Agent Requested:** Activated dynamically based on the semantic match of the rule's `description` to the query.
- **Manual:** Selected manually by the user in the UI.

### 5. Legacy rules file
- **Location:** `.cursorrules` at the project root.
- **Support:** Provides basic rule injection for legacy Cursor chat/features, but not used in Agent mode.
- **Action:** Included in scan for visibility (scan only), but never generated.

### 6. Ignore Rules
- **Location:** `.cursorignore` at the project root.
- **Format:** Standard gitignore syntax.
- **Action:** Excludes files and folders from being indexed, read, or explored by the Cursor agent.

---

## Phase 1 Implementation

In Phase 1, only scanning/parsing is implemented.
- `FindFiles` collects `.cursor/rules/*.mdc` files and `.cursorrules` (if present).
- `ParseRules` strips the YAML frontmatter from `.mdc` files and returns the rules as a `schema.Section`.
- Render is stubbed.

---

## Sources
- Cursor official documentation on Custom Rules (`.mdc` files).
- Cursor official documentation on Model Context Protocol (MCP) configuration.
- Community guidelines and GitHub structures for `.cursorignore` and `.cursorrules`.
