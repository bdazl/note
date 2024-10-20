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

	"github.com/spf13/cobra"
)

func noteMove(cmd *cobra.Command, args []string) {
	space, ids, err := checkMove(args)
	if err != nil {
		quitError("args", err)
	}

	db := dbOpen()
	defer db.Close()

	uniqueIds := removeDuplicates(ids)
	if err = db.MoveNotes(uniqueIds, space); err != nil {
		quitError("db move", err)
	}

	noteStr := "note"
	if len(uniqueIds) > 1 {
		noteStr = "notes"
	}
	fmt.Printf("Modified %s.\n", noteStr)
}

func checkMove(args []string) (string, []int, error) {
	if len(args) < 2 {
		return "", nil, fmt.Errorf("requires positional arguments space and id")
	}
	ids, err := parseIds(args[1:])
	if err != nil {
		return "", nil, err
	}

	return args[0], ids, nil
}
