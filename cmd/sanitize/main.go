// ABOUTME: CLI entry point for the sanitize tool. Parses subcommands and flags,
// reads stdin line-by-line, applies transformations, writes to stdout.
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tigger04/british-english-oed-fix/data"
	"github.com/tigger04/british-english-oed-fix/pkg/spelling"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: sanitize <subcommand> [flags]")
		fmt.Fprintln(os.Stderr, "subcommands: oed, symbols")
		os.Exit(2)
	}

	var doOED bool
	var quiet bool

	for _, arg := range os.Args[1:] {
		switch arg {
		case "oed":
			doOED = true
		case "-q":
			quiet = true
		default:
			fmt.Fprintf(os.Stderr, "unknown argument: %s\n", arg)
			os.Exit(2)
		}
	}

	if !doOED {
		fmt.Fprintln(os.Stderr, "usage: sanitize <subcommand> [flags]")
		fmt.Fprintln(os.Stderr, "subcommands: oed, symbols")
		os.Exit(2)
	}

	engine, err := spelling.NewOEDEngine(data.UsToUkData, data.IseToIzeData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for scanner.Scan() {
		line := scanner.Text()
		result := engine.ProcessLine(line)
		fmt.Fprintln(writer, result)
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		os.Exit(1)
	}

	if !quiet && engine.Changes > 0 {
		fmt.Fprintf(os.Stderr, "%d spelling corrections\n", engine.Changes)
	}
}
