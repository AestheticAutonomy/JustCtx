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

### 5. Background and Load Order
- `AGENTS.md` is an open, community-driven standard meant to act as a tool-agnostic "README for agents", reducing developer vendor lock-in.
- When co-existing with vendor-specific instruction files (like `CLAUDE.md` or `.cursorrules`), AI coding agents merge and prioritize rules, often combining them or reading `AGENTS.md` alongside native files.
- Some tools support hierarchical directories (scanning `AGENTS.md` at root and in subdirectories).

---

## Phase 1 Implementation

In Phase 1, only scanning/parsing is implemented.
- `FindFiles` locates `AGENTS.md` at the project root directory.
- `ParseRules` reads the rules as a raw `schema.Section`.
- Render is stubbed.

---

## Sources
- The open `AGENTS.md` specification and community guidelines (`https://agents.md`).
- Multi-agent tool-agnostic configuration standards.
