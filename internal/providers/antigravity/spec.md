# Antigravity Provider Specification

## Overview
The `antigravity` provider supports the Google Antigravity CLI (`agy`), which succeeded the legacy Gemini CLI on June 18, 2026. Antigravity CLI continues to read `GEMINI.md` as its primary project context file but introduces support for the cross-tool `AGENTS.md` standard. Additionally, it shifts workspace-specific configurations into the `.agents/` directory structure.

---

## File Conventions

### 1. Command Line Tool
*   **Binary:** `agy` (replacing the legacy `gemini` command).

### 2. Workspace Context & Rules Files
*   **Primary Context File:** `GEMINI.md` at the project root.
*   **Cross-Tool Context File:** `AGENTS.md` at the project root (also consumed by other AI coding assistants).
*   **Ignored Files:** `.antigravity.md` is not a standard convention and is ignored.

### 3. Project-Scoped Customizations (`.agents/` Directory)
Antigravity stores workspace-specific extensions under the `.agents/` folder:
*   **Workspace Skills:** `.agents/skills/<skill_name>/SKILL.md` — Custom workspace skill instructions and resources.
*   **Workspace MCP Config:** `.agents/mcp_config.json` — Workspace-scoped Model Context Protocol server definitions.

### 4. Global Configuration Paths
*   **Global Rules:** `~/.gemini/GEMINI.md` (Windows: `%USERPROFILE%\.gemini\GEMINI.md`).
*   **Global Customization Rules:** `~/.gemini/config/AGENTS.md` (Windows: `%USERPROFILE%\.gemini\config\AGENTS.md`).
*   **Global Settings:** `~/.gemini/antigravity-cli/settings.json` (Windows: `%USERPROFILE%\.gemini\antigravity-cli\settings.json`).
*   **Global MCP Config:** `~/.gemini/config/mcp_config.json` (Windows: `%USERPROFILE%\.gemini\config\mcp_config.json`).
*   **Global Skills:** `~/.gemini/config/skills/` (Windows: `%USERPROFILE%\.gemini\config\skills\`).

---

## Load Order & Precedence
When resolving instructions, the Antigravity CLI merges context from multiple files in the following order:
1.  **Global Scope:** Loads global rules from `~/.gemini/GEMINI.md` and global customizations from `~/.gemini/config/AGENTS.md`.
2.  **Project Root Scope:** Loads `[project-root]/GEMINI.md` and `[project-root]/AGENTS.md`. These append to/override global instructions.
3.  **Conflict/Priority:** When both `GEMINI.md` and `AGENTS.md` are present in the same directory, `GEMINI.md` takes precedence for conflicting directives. `AGENTS.md` is treated as supplementary cross-tool context and is appended after the native file's content.

---

## Migration from Gemini CLI
Developers migrating from Gemini CLI to Antigravity CLI follow these changes:
*   **Binary:** Command changes from `gemini <args>` to `agy <args>`.
*   **Workspace Folder:** Workspace-specific skills and configurations move from `.gemini/` to `.agents/` (e.g., rename `.gemini/skills/` to `.agents/skills/`).
*   **MCP Settings:** Workspace-scoped MCP configurations, which were historically nested inside `settings.json`, are moved into a separate `.agents/mcp_config.json` file.
*   **Automation:** The command `agy plugin import gemini` can be run to migrate legacy extensions and configurations.

---

## justctx Treatment & Phase 1 Implementation

`justctx` treats `antigravity` as a first-class provider for parsing and generating rules.
In Phase 1, only rules scanning/parsing is implemented.
*   **`FindFiles`:** Gathers:
    1.  Global rules: `~/.gemini/GEMINI.md`
    2.  Workspace rules: `[project-root]/GEMINI.md`
    3.  Workspace rules: `[project-root]/AGENTS.md`
*   **`ParseRules`:** Reads target files as raw `schema.Section` blocks.
*   **`RenderRules`:** Stubbed (returns `providers.ErrNotSupported`).

---

## Sources
*   *Google Developers Blog — Gemini CLI to Antigravity CLI transition:* [developers.googleblog.com](https://developers.googleblog.com/an-important-update-transitioning-gemini-cli-to-antigravity-cli/)
*   *Antigravity CLI GitHub repository:* [github.com/google-antigravity/antigravity-cli](https://github.com/google-antigravity/antigravity-cli)
*   *Antigravity CLI Documentation:* [antigravity.google/docs](https://antigravity.google/docs) — domain confirmed live; specific sub-paths (rules, mcp, plugins) exist per official site structure but content could not be fully retrieved at research time.
*   **Note:** Antigravity CLI launched June 18, 2026. All documentation reflects the official tooling as researched at that time.
