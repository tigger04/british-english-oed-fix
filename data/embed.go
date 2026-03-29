// ABOUTME: Embeds the word list files at compile time so the binary has no
// runtime file dependencies. Other packages import this for the raw data.
package data

import (
	_ "embed"
)

//go:embed us-to-uk.txt
var UsToUkData string

//go:embed ise-to-ize.txt
var IseToIzeData string
