//go:build !js && !wasm

package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/TheOpenDictionary/odict/lib/types"
	"github.com/TheOpenDictionary/odict/lib/utils"
	"github.com/fatih/color"
)

var faint = color.New(color.Faint)
var italic = color.New(color.Italic)
var italicFaint = color.New(color.Italic, color.Faint)
var italicFaintUnderlined = color.New(color.Italic, color.Faint, color.Underline)
var bold = color.New(color.Bold)
var boldUnderlined = color.New(color.Bold, color.Underline)
var parentheticalRegex = regexp.MustCompile(`(\(.*?\))`)

type PrintFormat = string

const (
	jsonFormat PrintFormat = "json"
	xmlFormat  PrintFormat = "xml"
	ppFormat   PrintFormat = "pp"
)

func ppExample(example string, underlined string, indent int) {
	start := strings.Index(strings.ToLower(example), strings.ToLower(underlined))

	fmt.Print(strings.Repeat(" ", indent))
	faint.Print("• ")

	if start >= 0 {
		end := start + len(underlined)
		italicFaint.Print(example[0:start])
		italicFaintUnderlined.Print(example[start:end])
		italicFaint.Printf("%s\n", example[end:])
	} else {
		italicFaint.Printf("%s\n", example)
	}

}

func ppDefinition(definition types.DefinitionRepresentable, numbering string, entry types.EntryRepresentable, indent int) {
	value := definition.Value
	matches := parentheticalRegex.FindAllStringIndex(value, -1)

	if len(matches) > 0 {
		j := 0

		fmt.Printf("%"+fmt.Sprint(indent)+"s. %s", numbering, value[0:matches[0][0]])

		for i := 0; i < len(matches); i += 1 {
			start := matches[i][0]
			end := matches[i][1]

			fmt.Print(value[j:start])
			italic.Print(value[start:end])

			j = end
		}

		fmt.Printf("%s\n", value[j:])
	} else {
		fmt.Printf("%"+fmt.Sprint(indent)+"s. %s\n", numbering, value)
	}

	for _, example := range definition.Examples {
		ppExample(example, entry.Term, indent+2)
	}
}

func ppGroup(group types.GroupRepresentable, i int, entry types.EntryRepresentable) {
	fmt.Printf("%4d. %s\n", i+1, group.Description)

	for j, definition := range group.Definitions {
		ppDefinition(definition, string(rune('a'+j)), entry, 7)
	}
}

func ppUsage(usage types.UsageRepresentable, entry types.EntryRepresentable) {
	italic.Printf("\n%s\n\n", usage.POS.Label)

	var i = 0

	for i < len(usage.Groups) {
		ppGroup(usage.Groups[i], i, entry)
		i++
	}

	for i < len(usage.Definitions) {
		ppDefinition(usage.Definitions[i], strconv.Itoa(i+1), entry, 4)
		i++
	}
}

func ppEty(ety types.EtymologyRepresentable, i int, showTitle bool, entry types.EntryRepresentable) {
	if showTitle {
		boldUnderlined.Printf("\nEtymology #%d\n", i+1)
	}

	if len(ety.Description) > 0 {
		fmt.Println(ety.Description)
	}

	for _, usage := range ety.Usages {
		ppUsage(usage, entry)
	}
}

func ppEntry(entry types.EntryRepresentable) {
	line := strings.Repeat("─", 32)

	fmt.Println(line)
	bold.Println(entry.Term)
	fmt.Println(line)

	etys := entry.Etymologies

	for i, ety := range etys {
		ppEty(ety, i, len(etys) > 1, entry)
	}

}

func prettyPrint(entries [][]types.EntryRepresentable) bool {
	fmt.Println()

	hasEntries := false

	for _, entry := range entries {
		for _, subentry := range entry {
			hasEntries = true
			ppEntry(subentry)
		}
	}

	return hasEntries
}

func PrintEntries(entries [][]types.EntryRepresentable, format PrintFormat, indent bool) {
	switch {
	case format == "json":
		fmt.Print(utils.SerializeToJSON(entries, indent))
	case format == "xml":
		fmt.Print(utils.SerializeToXML(entries, indent))
	case format == "pp":
		if !prettyPrint(entries) {
			fmt.Printf("No entries found!\n")
		}
	}
}
