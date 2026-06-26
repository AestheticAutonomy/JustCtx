# Claude Provider Specification

## File Conventions

### 1. Rules
- **Global Config:** `~/.claude/CLAUDE.md` (Windows equivalent: `%USERPROFILE%\.claude\CLAUDE.md`). Applied globally to all projects.
- **Project-Scoped:** `CLAUDE.md` or `.claude/CLAUDE.md` at the project root. Shared with the team via version control.
- **Local Overrides:** `CLAUDE.local.md` at the project root. Personal, git-ignored local overrides.
- **Subdirectory Rules:** `CLAUDE.md` files placed in subdirectories.
- **Loading Hierarchy & Lazy Loading:**
  - At session startup, Claude Code loads `~/.claude/CLAUDE.md`, the root project rules (`CLAUDE.md` / `.claude/CLAUDE.md`), `CLAUDE.local.md`, and any subdirectory `CLAUDE.md` files located in the ancestor folders of the current working directory.
  - Subdirectory `CLAUDE.md` files located in descendant folders (subdirectories) are **lazy loaded** on demand only when Claude’s tools (such as file reads/writes) navigate into or access files within those directories.
  - A directory hierarchy walk shows what *would* load for a given path, but not necessarily a guaranteed snapshot of a live session if those subdirectories haven't been accessed yet.

### 2. Model Context Protocol (MCP)
- **Project-Scoped:** `.mcp.json` at the project root. Can be committed to version control.
- **User-Scoped:** `~/.claude.json` (Windows: `%USERPROFILE%\.claude.json`), with servers defined under the `mcpServers` field. Available across all projects.
- **Note:** `.claude/mcp.json` is not supported/recognized and is silently ignored.

### 3. Hooks & Permissions
- **Settings Files:** Configured across four hierarchical scopes:
  1. **Managed:** Highest priority settings.
  2. **Local Settings:** `.claude/settings.local.json` (local project-scoped, git-ignored).
  3. **Project Settings:** `.claude/settings.json` (project-scoped, team-shared).
  4. **User Settings:** `~/.claude/settings.json` (global defaults).
- **Permissions:** Hard deny/allow/ask patterns are configured under the `permissions` field (e.g. `permissions.deny`).
  - Evaluated in the order: **Deny → Ask → Allow**. A deny rule always overrides allow or ask rules.
  - Supports wildcards, e.g. `Read(./secrets/**)` or `Bash(curl *)`.
  - Array-based settings like permissions merge across configuration levels (concatenated and deduplicated).
- **Hooks:** Automates shell actions at lifecycle events:
  - **`PreToolUse`:** Runs before a tool is executed. Can inspect input passed via JSON on `stdin` and block the operation by returning a non-zero exit code.
  - **`PostToolUse`:** Runs after a tool successfully completes.
  - Hook configurations specify: `matcher` (e.g., `Bash`, `Edit`, `*`), `type` (e.g., `command`), `command`, `timeout`, and `statusMessage`.
- **Exclusions (`.claudeignore`):**
  - A file named `.claudeignore` **is** natively supported at the project root for context management (excluding files from automatic discovery, indexing, or search).
  - However, it does not act as a security sandbox (files can still be read if explicitly requested). For strict read prevention, `permissions.deny` in `settings.json` or custom `PreToolUse` hook scripts must be used.

### 4. Custom Commands / Skills
- **Legacy custom commands:** Live in the `.claude/commands/` folder (or user-scoped at `~/.claude/commands/`) as Markdown files. Placeholders like `$ARGUMENTS`, `$1`, and shell injection (e.g., `! git diff`) are supported.
- **Skills System (Recommended):** Custom commands have been merged into the Skills system. Custom skills reside under `.claude/skills/<name>/` (project-scoped) or `~/.claude/skills/<name>/` (user-scoped).
  - Each skill directory contains a `SKILL.md` file (which contains frontmatter for metadata/matching rules and the instructions prompt) and optional supporting scripts/examples/resources.
  - Skills can be automatically loaded based on context matching of their description, in addition to being invocable via their name as a slash command.

---

## Phase 1 Implementation

In Phase 1, only the `rules` type (`CLAUDE.md` and subdirectory `CLAUDE.md` files) is implemented. All other types are stubbed out.

---

## Sources
- Anthropic Claude Code official documentation (Configuring Claude, Custom Commands, and MCP setup).
- Anthropic `claude-code` GitHub repository discussions & issues regarding `.claudeignore` security behavior and loader logic.
- JSON Schema Store definitions for Claude Code settings (`https://json.schemastore.org/claude-code-settings.json`).
