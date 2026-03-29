<!-- Version: 1.0 | Last updated: 2026-03-29 -->

# Vision

## Purpose

`sanitize` is a fast, composable CLI tool for converting English text to Oxford (OED) spelling. It exists because spell checkers are inconsistent about OED spelling, and most are unaware that `-ize` is the OED-preferred suffix in British English.

## Goals

1. **Correct, not clever.** The tool does dictionary lookups against curated word lists. It does not guess, infer, or apply heuristic rules. If a word isn't in the list, it passes through unchanged.

2. **Fast enough to be invisible.** The tool must be fast enough to use in interactive workflows (clipboard paste, editor integration) without perceptible delay. Target: <10ms for typical inputs.

3. **Composable.** stdin/stdout, Unix filter convention. Works in pipelines, editor integrations, git hooks, clipboard managers, and any other context that can pipe text.

4. **Two concerns, clearly separated.** Spelling correction (`oed`) and symbol sanitization (`symbols`) are independent subcommands that can be used alone or together. Neither assumes the other.

5. **Single source of truth.** The word lists in `data/` are the authoritative dictionaries. All consumers (the CLI tool, Hammerspoon, future integrations) read from these files. No duplication.

6. **Easy to maintain.** Adding a word means adding one line to a text file. No code changes required. The word lists are human-readable and diffable.

## Non-goals

- This is not a spell checker. It does not detect or highlight errors.
- This is not a grammar checker.
- This is not exhaustive. The dictionaries cover common words and will grow over time.
- This does not handle locale-specific vocabulary differences (e.g. "pavement" vs "sidewalk") — only spelling differences.

## Target audience

- Writers who follow OED/Oxford spelling and are tired of fighting spell checkers
- Programmers who want to sanitize typographic characters (smart quotes, em dashes) from text before it enters code or plain-text contexts
- Anyone who pipes text through Unix tools

## Consumers

- **CLI**: direct shell usage, pipelines
- **Hammerspoon**: clipboard paste integration via `hs.task.new()`
- **Future**: editor plugins, git hooks, CI linting
