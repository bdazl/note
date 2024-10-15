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
	PlainStyle Style = "plain"
	TitleStyle Style = "title"

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

	notes, err := collectNotes(args)
	if err != nil {
		quitError("collect notes", err)
	}

	pprintNotes(notes, style, color)
}

func styleColorOpts() (Style, bool, error) {
	var (
		style   Style
		doColor bool
	)
	switch styleArg {
	case string(PlainStyle):
		style = PlainStyle
	case string(TitleStyle):
		style = TitleStyle
	default:
		return "", false, fmt.Errorf("unrecognized style")
	}

	switch colorArg {
	case "auto":
		doColor = !color.NoColor
	case "yes", "always":
		doColor = true
	case "no", "never":
		doColor = false
	}

	return style, doColor, nil
}

func collectNotes(spaces []string) (db.Notes, error) {
	sortOpts, pageOpts, err := listOpts()
	if err != nil {
		return nil, fmt.Errorf("args: %w", err)
	}

	d := dbOpen()
	defer d.Close()

	lsSpaces := spaces
	if len(spaces) == 0 {
		allSpaces, err := d.ListSpaces(nil)
		if err != nil {
			return nil, fmt.Errorf("ls spaces: %w", err)
		}
		lsSpaces = allSpaces
	}

	notes, err := d.ListNotes(lsSpaces, sortOpts, pageOpts)
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
	case PlainStyle:
		printNotesPlain(notes)
	case TitleStyle:
		printNotesTitle(notes, doColor)
	}
}

func printNotesPlain(notes db.Notes) {
	for _, n := range notes {
		fmt.Printf(n.Content)
	}
}

func printNotesTitle(notes db.Notes, doColor bool) {
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

		// TODO: Why is a newline added?
		content := strings.TrimRight(n.Content, "\n")
		fmt.Println(content)
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
