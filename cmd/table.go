/*
Copyright © 2024 Jacob Peyron <jacob@peyron.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

const (
	Pin  = "📌"
	BoxH = "─"
)

var (
	idCol      = "ID"
	spaceCol   = "Space"
	pinCol     = "Pin"
	createdCol = "Created"
	// lastUpdatedCol = "Updated"
	previewCol = "Preview"

	tableCols = []string{idCol, spaceCol, pinCol, createdCol, previewCol}
)

func noteTable(cmd *cobra.Command, args []string) {
	notes, err := collectNotes(args)
	if err != nil {
		quitError("collect notes", err)
	}

	printTable(notes)
}

func printTable(notes db.Notes) {
	var (
		tw      = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		dateFmt = "2006-01-02"
	)

	// Print table header
	for _, col := range tableCols {
		fmt.Fprintf(tw, "%v\t", col)
	}

	fmt.Fprintln(tw)

	// Print box outline
	for _, col := range tableCols {
		line := strings.Repeat(string(BoxH), len(col))
		fmt.Fprintf(tw, "%s\t", line)
	}

	fmt.Fprintln(tw)

	// Print notes
	for _, note := range notes {
		pin := "no"
		if note.Pinned {
			pin = "yes" // Use Pin here, when a unicode-aware tabwriter is implemented
		}
		preview := getPreview(note.Content, int(previewArg))
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\n",
			note.ID, note.Space, pin,
			note.Created.Format(dateFmt),
			preview,
		)
	}

	tw.Flush()
}

func getPreview(content string, wordCount int) string {
	fields := strings.Fields(content)
	count := len(fields)
	if count < wordCount {
		return strings.Join(fields[:count], " ")
	}
	return strings.Join(fields[:wordCount], " ")
}
