# Project: british-english-oed-fix (sanitize)

## What this is

A Go CLI tool called `sanitize` that converts English text to OED (Oxford) spelling. Two subcommands: `oed` (spelling) and `symbols` (typographic character cleanup). stdin/stdout filter.

## Key context

- The word lists in `data/` are the single source of truth. No word pairs are hardcoded in Go source.
- Word lists use `wrong=correct` format, one pair per line. `#` for comments.
- The binary name is `sanitize`, not `british-english-oed-fix` (repo name differs from binary name).
- This project was extracted from a Hammerspoon config. The hammerspoon repo (tigger04/hammerspoon) has issue #26 tracking the migration to use this binary.

## Architecture

- `cmd/sanitize/main.go` — CLI entry point
- `pkg/spelling/` — spelling and symbol engines
- `data/us-to-uk.txt` — 721 US→UK spelling pairs
- `data/ise-to-ize.txt` — 1,213 -ise→-ize pairs
- Word lists embedded via `go:embed` — no runtime file dependencies

See `docs/architecture.md` for full details.

## Building and testing

```bash
make build    # compile binary
make test     # run regression tests
make install  # install locally
make release  # version bump + Homebrew formula update
```

## CLI design

```bash
sanitize oed              # spelling only
sanitize symbols          # symbol cleanup only
sanitize oed symbols      # both (any order)
sanitize                  # error: no subcommand
sanitize oed -q           # quiet (no stderr summary)
```

## Word list maintenance

To add a word:
1. Add a line to the appropriate file in `data/`
2. `make build` — the word list is embedded at compile time
3. Test and commit

## OED spelling note

OED (Oxford English Dictionary) prefers -ize endings (organize, recognize, standardize). This is NOT American spelling — it is the original British English standard, still used in technical writing, scientific publishing, and by international organizations. The tool converts non-OED -ise spellings to -ize, and US spellings (center, analyze) to UK (centre, analyse).
