// ABOUTME: Unit tests for the OED spelling engine.
// Tests word list loading, case-preserving replacement, and false positive avoidance.
package spelling

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// loadTestWordLists reads the word list files from the data/ directory.
func loadTestWordLists(t *testing.T) (*OEDEngine, error) {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	usToUk, err := os.ReadFile(filepath.Join(repoRoot, "data", "us-to-uk.txt"))
	if err != nil {
		t.Fatalf("failed to read us-to-uk.txt: %v", err)
	}
	iseToIze, err := os.ReadFile(filepath.Join(repoRoot, "data", "ise-to-ize.txt"))
	if err != nil {
		t.Fatalf("failed to read ise-to-ize.txt: %v", err)
	}

	return NewOEDEngine(string(usToUk), string(iseToIze))
}

// RT-001: ise-to-ize list — lowercase "organise" → "organize"
func TestOed_LowercaseIseWord_ConvertsToIze_RT001(t *testing.T) {
	engine, err := loadTestWordLists(t)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	result := engine.ProcessLine("organise")
	if result != "organize" {
		t.Errorf("got %q, want %q", result, "organize")
	}
}

// RT-005: us-to-uk list — lowercase "center" → "centre"
func TestOed_LowercaseUsWord_ConvertsToUk_RT005(t *testing.T) {
	engine, err := loadTestWordLists(t)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	result := engine.ProcessLine("center")
	if result != "centre" {
		t.Errorf("got %q, want %q", result, "centre")
	}
}

// RT-002: uppercase "ORGANISE" → "ORGANIZE"
func TestOed_UppercaseWord_PreservesCase_RT002(t *testing.T) {
	engine, err := loadTestWordLists(t)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	result := engine.ProcessLine("ORGANISE")
	if result != "ORGANIZE" {
		t.Errorf("got %q, want %q", result, "ORGANIZE")
	}
}

// RT-003: title case "Organise" → "Organize"
func TestOed_TitleCaseWord_PreservesCase_RT003(t *testing.T) {
	engine, err := loadTestWordLists(t)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	result := engine.ProcessLine("Organise")
	if result != "Organize" {
		t.Errorf("got %q, want %q", result, "Organize")
	}
}

// RT-004: words containing "-ise" that are NOT in the list remain unchanged
func TestOed_WordsWithIseSubstringNotInList_Unchanged_RT004(t *testing.T) {
	engine, err := loadTestWordLists(t)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	input := "advise promise enterprise"
	result := engine.ProcessLine(input)
	if result != input {
		t.Errorf("got %q, want %q", result, input)
	}
}

// RT-009: lorem ipsum passthrough — no modifications to unrelated text
func TestOed_LoremIpsum_PassesThrough_RT009(t *testing.T) {
	engine, err := loadTestWordLists(t)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	input := "Lorem ipsum dolor sit amet"
	result := engine.ProcessLine(input)
	if result != input {
		t.Errorf("got %q, want %q", result, input)
	}
}
