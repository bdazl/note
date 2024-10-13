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
	"github.com/spf13/viper"
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
)

func noteList(cmd *cobra.Command, args []string) {
	style, color, err := styleColorOpts()
	if err != nil {
		quitError("args", err)
	}

	notes, err := collectNotes()
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

func collectNotes() ([]db.Note, error) {
	sortOpts, pageOpts, err := listOpts()
	if err != nil {
		return nil, fmt.Errorf("args: %w", err)
	}

	d := dbOpen()
	defer d.Close()

	spaces, err := getSpaces(d)
	if err != nil {
		return nil, fmt.Errorf("ls spaces: %w", err)
	}

	notes, err := d.ListNotes(spaces, sortOpts, pageOpts)
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

func pprintNotes(notes []db.Note, style Style, doColor bool) {
	switch style {
	case PlainStyle:
		printNotesPlain(notes)
	case TitleStyle:
		printNotesTitle(notes, doColor)
	}
}

func printNotesPlain(notes []db.Note) {
	for _, n := range notes {
		fmt.Printf(n.Content)
	}
}

func printNotesTitle(notes []db.Note, doColor bool) {
	var (
		col = color.New(color.FgGreen)

		boxThin = '\u2500'
		line    = strings.Repeat(string(boxThin), 4)
	)

	// fatih/color is (too) smart and disables colors for non-terminal outputs
	// doColor considers such cases when color is set to 'auto'. If doColor is
	// true - we should explicitly enable the color
	if doColor {
		color.NoColor = false
	}

	for _, n := range notes {
		if doColor {
			col.Printf("%s [", line)
			fmt.Printf("%v", n.ID)
			col.Printf("] %s\n", line)
		} else {
			fmt.Printf("%s [%v] %s\n", line, n.ID, line)
		}

		// TODO: Why is a newline added?
		content := strings.TrimRight(n.Content, "\n")
		fmt.Println(content)
	}

	if doColor {
		color.NoColor = true
	}
}

func getSpaces(d *db.DB) ([]string, error) {
	if allArg {
		spaces, err := d.ListSpaces(nil)
		if err != nil {
			return nil, err
		}
		return spaces, nil
	}
	return viper.GetStringSlice(ViperListSpaces), nil
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
