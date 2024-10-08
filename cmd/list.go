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
	"golang.org/x/exp/maps"
)

var (
	validSortColumns = map[string]db.NoteColumn{
		"id":       db.ColumnID,
		"creation": db.ColumnCreatedAt,
		"updated":  db.ColumnCreatedAt,
		"title":    db.ColumnTitle,
		"archived": db.ColumnIsArchived,
		"favorite": db.ColumnIsFavorite,
	}
)

func noteList(cmd *cobra.Command, args []string) {
	sortColumn, err := checkSortArguments()
	if err != nil {
		quitError("args", err)
	}

	d, err := db.Open(storagePath)
	if err != nil {
		quitError("db open", err)
	}

	notes, err := db.ListNotes(d, sortColumn, !descending, limit, offset)
	if err != nil {
		quitError("db list", err)
	}

	pprintNotes(notes)
}

func checkSortArguments() (db.NoteColumn, error) {
	sortColumn, err := mapSortBy(sortBy)
	if err != nil {
		return "", err
	}

	if limit < 0 {
		return "", fmt.Errorf("limit must be zero or positive")
	} else if limit == 0 && offset != 0 {
		return "", fmt.Errorf("offset is only valid if you specify a limit")
	} else if limit > 0 && offset < 0 {
		return "", fmt.Errorf("offset must be zero or positive")
	}

	return sortColumn, nil
}

func pprintNotes(notes []db.Note) {
	for _, n := range notes {
		if n.Title != nil && *n.Title != "" {
			fmt.Printf("%v: ", *n.Title)
		}
		fmt.Println(n.Note)
	}
}

func mapSortBy(s string) (db.NoteColumn, error) {
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
