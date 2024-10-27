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
	"strings"

	"github.com/bdazl/note/db"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type Color string
type Style string

const (
	RawStyle   Style = "raw"
	LightStyle Style = "light"
	FullStyle  Style = "full"

	AutoColor   Color = "auto"
	NeverColor  Color = "never"
	AlwaysColor Color = "always"
)

var (
	validNoteSortColumns = map[string]db.Column{
		"id":      db.IDColumn,
		"space":   db.SpaceColumn,
		"content": db.ContentColumn,
		"created": db.CreatedColumn,
		"updated": db.CreatedColumn,
	}

	Green = color.New(color.FgGreen)
)

func noteList(cmd *cobra.Command, args []string) {
	style, color, err := styleColorOpts()
	if err != nil {
		quitError("args", err)
	}

	notes, err := selectNotes(args)
	if err != nil {
		quitError("collect notes", err)
	}

	pprintNotes(notes, style, color)
}

func selectNotes(spaces []string) (db.Notes, error) {
	sortOpts, pageOpts, err := listOpts()
	if err != nil {
		return nil, fmt.Errorf("args: %w", err)
	}

	d := dbOpen()
	defer d.Close()

	notes, err := d.SelectNotes(spaces, allArg, sortOpts, pageOpts)
	if err != nil {
		return nil, fmt.Errorf("db list: %w", err)
	}
	return notes, nil
}

func listOpts() (*db.SortOpts, *db.PageOpts, error) {
	sortColumn, err := mapNoteSortColumn(sortByArg)
	if err != nil {
		return nil, nil, err
	}
	sortOpts := &db.SortOpts{
		Ascending:  !descendingArg,
		SortColumn: sortColumn,
	}
	if err = sortOpts.Check(); err != nil {
		return nil, nil, err
	}
	pageOpts := &db.PageOpts{
		Limit:  limitArg,
		Offset: offsetArg,
	}
	if err = pageOpts.Check(); err != nil {
		return nil, nil, err
	}
	return sortOpts, pageOpts, nil
}

func pprintNotes(notes db.Notes, style Style, doColor bool) {
	switch style {
	case RawStyle:
		printNotesRaw(notes)
	case LightStyle:
		printNotesLight(notes, doColor)
	case FullStyle:
		printNotesFull(notes, doColor)
	}
}

func printNotesRaw(notes db.Notes) {
	for _, n := range notes {
		content := strings.TrimRight(n.Content, "\n")
		fmt.Println(content)
	}
}

func printNotesLight(notes db.Notes, doColor bool) {
	var (
		line = strings.Repeat(BoxH, 4)
	)

	preColor(doColor)

	for _, n := range notes {
		pinned := ""
		if n.Pinned {
			pinned = Pin
		}
		if doColor {
			Green.Printf("%s [", line)
			fmt.Printf("%v", n.ID)
			Green.Printf("] %s", line)
			fmt.Printf(" %s\n", pinned)
		} else {
			fmt.Printf("%s [%v] %s\n", line, n.ID, line)
		}

		content := strings.TrimRight(n.Content, "\n")
		fmt.Println(content)
	}

	postColor(doColor)
}

func printNotesFull(notes db.Notes, doColor bool) {
	fullFmt := "2006-01-02 15:04:05"
	preColor(doColor)

	for _, n := range notes {
		created := n.Created.Format(fullFmt)
		updated := n.LastUpdated.Format(fullFmt)
		pinned := "no"
		if n.Pinned {
			pinned = "yes"
		}
		if doColor {
			Green.Printf("ID: ")
			fmt.Printf("%v\n", n.ID)
			Green.Printf("Pinned: ")
			fmt.Printf("%v\n", pinned)
			Green.Printf("Space: ")
			fmt.Printf("%v\n", n.Space)
			Green.Printf("Created: ")
			fmt.Printf("%v\n", created)
			Green.Printf("Last Updated: ")
			fmt.Printf("%v\n", updated)
			Green.Printf("Content:\n")
			fmt.Printf("%v\n", n.Content)
		} else {
			fmt.Printf(
				"ID: %v\nPinned: %v\nSpace: %v\nCreated: %v\nLast Updated: %v\nContent:\n%v\n",
				n.ID, pinned, n.Space, created, updated, n.Content)
		}
	}

	postColor(doColor)
}

func preColor(doColor bool) {
	// fatih/color is (too) smart and disables colors for non-terminal outputs
	// doColor considers such cases when color is set to 'auto'. If doColor is
	// true - we should explicitly enable the color
	if doColor {
		color.NoColor = false
	}
}

func postColor(doColor bool) {
	if doColor {
		color.NoColor = true
	}
}

func mapNoteSortColumn(s string) (db.Column, error) {
	out, ok := validNoteSortColumns[s]
	if !ok {
		return "", fmt.Errorf("invalid sort option: %v", s)
	}
	return out, nil
}

func getSortKeys() string {
	sortKeys := maps.Keys(validNoteSortColumns)
	return strings.Join(sortKeys, ", ")
}
