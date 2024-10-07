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

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

func noteRoot(cmd *cobra.Command, args []string) {
	pprintNotes(list(10, 0))
}

func noteList(cmd *cobra.Command, args []string) {
	pprintNotes(list(0, 0))
}

func pprintNotes(notes []db.Note) {
	for _, n := range notes {
		if n.Title != nil && *n.Title != "" {
			fmt.Printf("%v: ", *n.Title)
		}
		fmt.Println(n.Note)
	}
}

func list(limit, offset int) []db.Note {
	d, err := db.Open(storagePath)
	if err != nil {
		quitError("db open", err)
	}

	notes, err := db.ListNotes(d, db.ColumnCreatedAt, false, limit, offset)
	if err != nil {
		quitError("db list", err)
	}
	return notes
}
