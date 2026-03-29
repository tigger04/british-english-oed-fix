<!-- Version: 1.0 | Last updated: 2026-03-29 -->

# Testing Strategy

## Framework

Go's built-in `testing` package. No external test framework.

```bash
make test          # run regression tests
make test-one-off  # run one-off tests
```

## Test types

### Unit tests

Each package has `_test.go` files alongside the source:

- `pkg/spelling/oed_test.go` — word list loading, map building, case-preserving replacement
- `pkg/spelling/symbols_test.go` — character replacement
- `pkg/spelling/engine_test.go` — orchestration, combined mode

### Integration tests

End-to-end tests that invoke the compiled binary:

- `tests/regression/` — stdin/stdout pipeline tests, flag parsing, error cases
- Tests run the actual binary with known inputs and compare outputs

### Data quality tests

Verify the word lists themselves:

- No duplicate keys
- Every line matches expected format
- Spot-check known conversions
- No known-bad entries

## Test structure

Follow Arrange-Act-Assert:

```go
func TestOedConvertsOrganiseToOrganize(t *testing.T) {
    // Arrange
    engine := spelling.NewEngine(spelling.WithOED())

    // Act
    result := engine.ProcessLine("I need to organise this")

    // Assert
    if result != "I need to organize this" {
        t.Errorf("got %q, want %q", result, "I need to organize this")
    }
}
```

## Test naming

```
Test<Unit>_<Scenario>_<Expected>
```

Examples:
- `TestOed_LowercaseIseWord_ConvertsToIze`
- `TestOed_WordNotInList_PassesThrough`
- `TestSymbols_SmartDoubleQuotes_ConvertToStraight`
- `TestCLI_NoSubcommand_ExitsWithError`

## Coverage targets

- Minimum 80% line coverage for new code
- 100% coverage for the replacement engine (it's the core logic)
- Word list data quality tests cover every known edge case from the curation analysis

## Test IDs

Tests carry IDs per project convention:

| Prefix | Type | Location | Run by |
|--------|------|----------|--------|
| `RT-NNN` | Regression | `tests/regression/` | `make test` |
| `OT-NNN` | One-off | `tests/one_off/` | `make test-one-off` |
| `UT-NNN` | User (manual) | issue AC table only | human |

ID allocation tracked in `tests/NEXT_IDS.txt`.

## False positive testing

Critical: words that share substrings with dictionary entries must not be affected. Test with:
- "advise" (contains "ise" but is not in the -ise→-ize list)
- "promise", "enterprise", "compromise" (same)
- "advertisement" (contains "ise" substring)
- Lorem ipsum text (no changes expected)
