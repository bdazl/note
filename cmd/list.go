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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var (
	validSortColumns = map[string]db.NoteColumn{
		"id":      db.ColumnID,
		"space":   db.ColumnSpace,
		"content": db.ColumnContent,
		"created": db.ColumnCreatedAt,
		"updated": db.ColumnCreatedAt,
	}
)

func noteList(cmd *cobra.Command, args []string) {
	sortOpts, pageOpts, err := listOpts()
	if err != nil {
		quitError("args", err)
	}

	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	spaces, err := getSpaces(d)
	if err != nil {
		quitError("ls spaces", err)
	}

	notes, err := d.ListNotes(spaces, sortOpts, pageOpts)
	if err != nil {
		quitError("db list", err)
	}

	pprintNotes(notes)
}

func noteSpaces(cmd *cobra.Command, args []string) {
	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	spaces, err := d.ListSpaces()
	if err != nil {
		quitError("db list", err)
	}

	spacesStr := strings.Join(spaces, " ")
	fmt.Println(spacesStr)
}

func listOpts() (*db.SortOpts, *db.PageOpts, error) {
	sortColumn, err := mapSortColumn(sortByArg)
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

func pprintNotes(notes []db.Note) {
	for _, n := range notes {
		fmt.Println(n.Content)
	}
}

func getSpaces(d *db.DB) ([]string, error) {
	if allArg {
		spaces, err := d.ListSpaces()
		if err != nil {
			return nil, err
		}
		return spaces, nil
	}
	return viper.GetStringSlice(ViperListSpaces), nil
}

func mapSortColumn(s string) (db.NoteColumn, error) {
	out, ok := validSortColumns[s]
	if !ok {
		return "", fmt.Errorf("invalid sort option: %v", s)
	}
	return out, nil
}

func getSortKeys() string {
	sortKeys := maps.Keys(validSortColumns)
	return strings.Join(sortKeys, ", ")
}
