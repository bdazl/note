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

	"github.com/spf13/cobra"
)

func noteClean(cmd *cobra.Command, args []string) {
	db := dbOpen()
	defer db.Close()

	// allInSpaceArg is set and no args provided means find all notes in specific space
	allNotesInSpace, err := db.SelectNotes([]string{TrashSpace}, true, nil, nil)
	if err != nil {
		quitError("db list", err)
	}

	ids := allNotesInSpace.GetIDs()
	uniqueIds := removeDuplicates(ids)
	if len(uniqueIds) == 0 {
		fmt.Println("Trash is empty")
		os.Exit(0)
	}

	if !noConfirmArg {
		fmt.Printf("WARNING: You are about to permanently remove %v note(s).\n", len(uniqueIds))
		fmt.Printf("Write 'yes' to confirm permanent delete: ")
		response := readUserInput()
		if response != "yes" {
			os.Exit(2)
		}
	}

	if err := db.PermanentRemoveNotes(uniqueIds); err != nil {
		quitError("db remove", err)
	}

	count := len(uniqueIds)
	if count == 1 {
		fmt.Printf("Note removed from %v\n", TrashSpace)
	} else {
		fmt.Printf("%v notes removed from %v\n", count, TrashSpace)
	}
}
