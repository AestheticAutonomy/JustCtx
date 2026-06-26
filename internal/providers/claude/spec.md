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
- Configured in `.claude/settings.json`.

### 4. Custom Commands
- Custom slash commands live in the `.claude/commands/` folder.

---

## Phase 1 Implementation

In Phase 1, only the `rules` type (`CLAUDE.md`) is implemented. All other types are stubbed out.
