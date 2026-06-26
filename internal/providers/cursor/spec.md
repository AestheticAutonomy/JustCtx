# Cursor Provider Specification

## File Conventions

### 1. Primary Format: `.cursor/rules/*.mdc`
- **Location:** `.cursor/rules/` under the project root.
- **Structure:** MDC rules are markdown files ending in `.mdc` with a YAML frontmatter block containing metadata:
  ```yaml
  ---
  description: "Rule description used by the LLM for request-based activation"
  globs: "glob patterns (e.g. *.go, internal/**/*.go)"
  alwaysApply: false
  ---
  ```
- **Plain Markdown:** Plain `.md` files in `.cursor/rules/` are silently ignored in Agent mode (frontmatter is required).

### 2. Activation Modes
Cursor uses four activation modes for MDC rules:
- **Always Apply:** `alwaysApply: true` in frontmatter.
- **Auto Attached:** File matches the pattern in `globs`.
- **Agent Requested:** Activated by the agent based on matching the rule's `description`.
- **Manual:** Selected manually by the user in the UI.

### 3. Legacy rules file
- **Location:** `.cursorrules` at the project root.
- **Support:** basic rule injection for legacy cursor chat, but not used in agent mode.
- **Action:** Included in scan for visibility (scan only), but never generated.

---

## Phase 1 Implementation

In Phase 1, only scanning/parsing is implemented.
- `FindFiles` collects `.cursor/rules/*.mdc` files and `.cursorrules` (if present).
- `ParseRules` strips the YAML frontmatter from `.mdc` files and returns the rules as a `schema.Section`.
- Render is stubbed.
