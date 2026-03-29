<!-- Version: 1.0 | Last updated: 2026-03-29 -->

# Implementation Plan

## Phases

### Phase 1: Core spelling engine

**Goal**: `sanitize oed` works — reads stdin, applies both word lists, writes corrected text to stdout.

- Parse and embed both word lists into a lookup map
- Implement word extraction and case-preserving replacement
- Implement line-by-line stdin→stdout processing
- Implement change summary on stderr
- Implement `-q` flag
- Tests: word list quality, case preservation, false positive avoidance, pipeline behaviour

### Phase 2: Symbol sanitization

**Goal**: `sanitize symbols` works — replaces typographic characters with ASCII equivalents.

- Implement character replacement table
- Implement bullet-to-hyphen line-start detection
- Tests: each character type, passthrough of unaffected text

### Phase 3: Combined mode and CLI polish

**Goal**: `sanitize oed symbols` works, `sanitize` with no args errors, `-h` and `--version` work.

- Implement subcommand parsing (any order, deduplication)
- Fixed transformation order (spelling then symbols)
- Error on no subcommands
- Help text and version flag
- Tests: combined mode, flag parsing, error cases

### Phase 4: Homebrew distribution

**Goal**: `brew install sanitize` works.

- Create Homebrew tap formula
- `make release` workflow (version bump, build, SHA256, formula update, tag)
- `make install` for local development
- Test: install from tap, run, verify output

### Phase 5: Hammerspoon integration (separate repo)

**Goal**: Hammerspoon's clipboard paste uses the installed `sanitize` binary.

- Tracked in hammerspoon repo issue #26
- Replace inline Lua spelling logic with `hs.task.new()` call
- Remove legacy spelling files from hammerspoon repo

## Dependencies

```
Phase 1 ← Phase 2 (can be parallel, no code dependency)
Phase 1 + 2 ← Phase 3 (needs both engines)
Phase 3 ← Phase 4 (needs working binary)
Phase 4 ← Phase 5 (needs Homebrew install)
```

Phases 1 and 2 can be developed in parallel. Phase 3 integrates them.

## Issue tracking

Each phase maps to one or more GitHub issues on this repo. Issues follow the standards in CLAUDE.md (AC table, TDD, solution outline).
