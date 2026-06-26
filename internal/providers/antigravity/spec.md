# Antigravity Provider Specification

## File Conventions

### 1. Binary
- `agy`

### 2. Workspace Files
- **Primary:** `GEMINI.md` at the project root.
- **Cross-Tool Standard:** `AGENTS.md` at the project root (also written/read by other providers).

### 3. Global Configuration
- **Global rules:** `~/.gemini/GEMINI.md` (Windows: `%USERPROFILE%\.gemini\GEMINI.md`) and the global customizations root `~/.gemini/config/AGENTS.md`.
- **Global Settings:** managed in `~/.gemini/antigravity-cli/settings.json`.

### 4. Skills (Phase 2+)
- **Location:** `.agents/skills/` (project-scoped) and `~/.gemini/config/skills/` (global).

### 5. MCP (Phase 2+)
- **Location:** `.agents/mcp_config.json`.

### 6. Ignored files
- `.antigravity.md` does not exist as a convention and is ignored.

### 7. Migration Notes
- Gemini CLI transitioned to Antigravity CLI as of June 18 2026.
- The underlying file formats and locations (`GEMINI.md` and `AGENTS.md`) carry over completely unchanged.
- Plugins and skills migration: `agy plugin import gemini` converts legacy extensions; old skill directories can be moved via `mv .gemini/skills/ .agents/skills/`.

---

## Phase 1 Implementation

Only the `rules` type (`GEMINI.md` and `AGENTS.md`) is implemented.
- `FindFiles` collects `GEMINI.md` and `AGENTS.md` from the project root and `~/.gemini/GEMINI.md` (global).
- `ParseRules` reads files as raw `schema.Section` blocks.
- Render is stubbed.

---

## Sources
- Antigravity documentation home (`https://antigravity.google/docs`).
- Customizations specifications (AGENTS.md workspace vs global configuration scopes).
- Antigravity CLI transition guides and command reference.
