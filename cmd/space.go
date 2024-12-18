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
	"sort"
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

func noteSpace(cmd *cobra.Command, args []string) {
	sortOpts, err := spacesSortOpt()
	if err != nil {
		quitError("arg", err)
	}

	// Parse ids, if argument(s) are supplied
	var findIds []int
	if len(args) > 0 {
		ids, err := parseIds(args)
		if err != nil {
			quitError("parse ids", err)
		}

		findIds = removeDuplicates(ids)
	}

	d := dbOpen()
	defer d.Close()

	// Either find spaces linked to notes
	var spaces []string
	if len(findIds) > 0 {
		notes, err := d.GetNotes(findIds)
		if err != nil {
			quitError("db get", err)
		}

		// When we get notes, the order is one-to-one, so we need to manually sort
		// and remove duplicates.
		allSpaces := notes.GetSpaces()
		spaces = removeDuplicates(allSpaces)
		if descendingArg {
			sort.Sort(sort.Reverse(sort.StringSlice(spaces)))
		} else {
			sort.Strings(spaces)
		}
	} else {
		// Or list all (or at least some) spaces
		lsSpaces, err := d.SelectSpaces(allArg, sortOpts)
		if err != nil {
			quitError("db list", err)
		}

		spaces = lsSpaces
	}

	// Print listed spaces
	if listArg {
		for _, space := range spaces {
			fmt.Println(space)
		}
	} else {
		spacesStr := strings.Join(spaces, " ")
		fmt.Println(spacesStr)
	}
}

func spacesSortOpt() (*db.SortOpts, error) {
	sortOpts := &db.SortOpts{
		Ascending:  !descendingArg,
		SortColumn: db.SpaceColumn,
	}

	// This should never error, because the column is hard-coded.
	// Let's be mindful of any future slip-ups :)
	if err := sortOpts.Check(); err != nil {
		return nil, err
	}

	return sortOpts, nil
}
