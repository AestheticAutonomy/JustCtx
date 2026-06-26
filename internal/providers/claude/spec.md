# Claude Provider Specification

## File Conventions

### 1. Rules
- **Project Root & Subdirectories:** `CLAUDE.md` files placed in the project root and subdirectories.
- **Global Config:** `~/.claude/CLAUDE.md` (Windows equivalent: `%USERPROFILE%\.claude\CLAUDE.md`).
- **Caveat:** Subdirectory `CLAUDE.md` files are loaded lazily by Claude (i.e. only when a file in that directory is active). A scan of the directory hierarchy shows what *would* load for a given path, but not necessarily a guaranteed snapshot of a live session.

### 2. Model Context Protocol (MCP)
- **Project-Scoped:** `.mcp.json` at the project root.
- **User-Scoped:** `~/.claude.json` (Windows: `%USERPROFILE%\.claude.json`), with servers defined under the `mcpServers` field.
- **Note:** `.claude/mcp.json` is not supported/recognized and is silently ignored.

### 3. Hooks & Permissions
- **Settings files:** Configured in `.claude/settings.json` (project-scoped, team shared) and `.claude/settings.local.json` (local project scoped, git-ignored). Global hooks and settings are in `~/.claude/settings.json`.
- **Permissions:** Hard deny patterns are configured in `settings.json` under `permissions.deny` (e.g. `Read(./build/)`).
- **Exclusions:** A file named `.claudeignore` is not natively supported out of the box by Claude Code itself; native file exclusion must be declared via the `permissions.deny` configuration or implemented using custom `PreToolUse` hooks (such as using the third-party `claudeignore` npm wrapper).

### 4. Custom Commands / Skills
- Custom slash commands live in the `.claude/commands/` folder or as system prompts inside `.claude/skills/<name>/SKILL.md`.

---

## Phase 1 Implementation

In Phase 1, only the `rules` type (`CLAUDE.md`) is implemented. All other types are stubbed out.

---

## Sources
- Claude Code official configuration documentation (e.g., config commands, custom commands, and MCP setup guides).
- Claude Code settings and hook structures verified against GitHub documentation and community packages (`claudeignore`).
