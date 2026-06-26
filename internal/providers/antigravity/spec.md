# Antigravity Provider Specification

## File Conventions

### 1. Binary
- `agy`

### 2. Workspace Files
- **Primary:** `GEMINI.md` at the project root.
- **Cross-Tool Standard:** `AGENTS.md` at the project root (also written/read by other providers).

### 3. Global Configuration
- **Global:** `~/.gemini/GEMINI.md` (Windows equivalent: `%USERPROFILE%\.gemini\GEMINI.md`).

### 4. Skills (Phase 2+)
- **Location:** `.agents/skills/`.

### 5. MCP (Phase 2+)
- **Location:** `.agents/mcp_config.json`.

### 6. Ignored files
- `.antigravity.md` does not exist as a convention and is ignored.

### 7. Migration Notes
- Gemini CLI transitioned to Antigravity CLI as of June 18 2026.
- The underlying file formats and locations (`GEMINI.md` and `AGENTS.md`) carry over completely unchanged.

---

## Phase 1 Implementation

Only the `rules` type (`GEMINI.md` and `AGENTS.md`) is implemented.
- `FindFiles` collects `GEMINI.md` and `AGENTS.md` from the project root and `~/.gemini/GEMINI.md` (global).
- `ParseRules` reads files as raw `schema.Section` blocks.
- Render is stubbed.
