# Agents Provider Specification

## Overview
`AGENTS.md` is an open, community-driven, tool-agnostic configuration standard designed to serve as a "README for AI agents". Originally introduced by OpenAI in August 2025 and subsequently transitioned to the Linux Foundation's Agentic AI Foundation, the standard aims to solve vendor-specific configuration fragmentation by providing a unified file for instructing AI coding assistants.

---

## File Conventions

### 1. File Location & Scope
*   **Project Root:** `[project-root]/AGENTS.md` — The canonical file for repository-scoped context.
*   **Subdirectories:** `[project-root]/path/to/subdir/AGENTS.md` — Folder-specific context (e.g., scoping different rules for `tests/` vs `src/`).
*   **Global Path (User-Level):**
    *   Linux/macOS: `~/.agents/AGENTS.md`
    *   Windows: `%USERPROFILE%\.agents\AGENTS.md`

### 2. Consumers & Official Support
Confirmed tools and assistants that read `AGENTS.md` natively or via standard configurations:
*   **Antigravity CLI (`agy`):** Natively reads project root and global configurations.
*   **Cursor (Agent mode):** Natively reads project root and merges it.
*   **GitHub Copilot:** Natively reads and parses project root context.
*   **Aider & OpenAI Codex:** Natively consume root files to seed system prompt instructions.
*   **Claude Code:** Does not natively discover `AGENTS.md` directly. Instead, compatibility is achieved by importing it at the top of `CLAUDE.md` using the `@AGENTS.md` directive or via symbolic linking.

### 3. Format
*   **Type:** Plain Markdown.
*   **Frontmatter:** None required or officially specified.
*   **Structure:** Standard Markdown headings (e.g., `# Tech Stack`, `## Coding Style`, `## Build Commands`) are parsed as logical section divisions by tools.

### 4. Load Order & Precedence
When co-existing with native vendor-specific files (e.g., `CLAUDE.md` or `.cursorrules` / `.cursor/rules/*.mdc`):
1.  **Global vs. Project:** Project-specific files (`[project-root]/AGENTS.md`) append to or override settings from the Global User-Level file (`~/.agents/AGENTS.md`).
2.  **Hierarchical Resolution:** Agents walk the directory tree from the root down to the current working directory, loading and merging all active `AGENTS.md` files (nearer files take precedence).
3.  **Vendor Precedence:** Vendor-specific instructions (like `CLAUDE.md` for Claude Code or `.cursorrules` for Cursor) typically take precedence over or are merged alongside the shared `AGENTS.md`.

---

## justctx Treatment

`justctx` treats the `agents` provider as a cross-tool, generic target for generation and a source for scanning.
*   **No Frills:** It does not introduce provider-specific custom directives or metadata flags.
*   **Scan Source:** Parsed to construct the consolidated context envelope.
*   **Generate Target:** Assembled guidelines are compiled into a standard Markdown format and output to the target `AGENTS.md` file.

---

## Phase 1 Implementation

In Phase 1, only scanning/parsing is implemented.
*   **`FindFiles`:** Locates the `AGENTS.md` at the project repository root.
*   **`ParseRules`:** Reads the rules from `AGENTS.md` as a raw `schema.Section` (no heading-based slicing is applied in Phase 1).
*   **`RenderRules`:** Stubbed (returns `providers.ErrNotSupported`).

---

## Sources
*   *Linux Foundation Agentic AI Foundation / Open Agentic Standard:* [agentsstandard.com](https://agentsstandard.com)
*   *OpenAI AGENTS.md Specification:* [agents.md](https://agents.md)
*   *Claude Code Import Convention:* Per `internal/providers/claude/spec.md` — Claude Code supports `@<path>` import directives in `CLAUDE.md` files; `@AGENTS.md` at the top of a root `CLAUDE.md` is the standard compatibility shim for cross-tool context.
*   *Cursor Rules Documentation:* [docs.cursor.com](https://docs.cursor.com)
