# Agents Provider Specification

## File Conventions

### 1. File Location
- **Location:** `AGENTS.md` at the project root.

### 2. Consumers
- Read by Claude Code, Cursor (Agent mode), and Antigravity CLI.

### 3. Cross-Tool Standard
- Not owned by any single provider.
- When generating files for multiple target tools, `AGENTS.md` is written only once (it contains cross-tool shared context and guidelines).

### 4. Format
- Plain markdown, no special frontmatter required.

---

## Phase 1 Implementation

In Phase 1, only scanning/parsing is implemented.
- `FindFiles` locates `AGENTS.md` at the project root directory.
- `ParseRules` reads the rules as a raw `schema.Section`.
- Render is stubbed.
