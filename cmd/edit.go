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
	"strconv"

	"github.com/spf13/cobra"
)

func noteEdit(cmd *cobra.Command, args []string) {
	id, err := checkEdit(args)
	if err != nil {
		quitError("args", err)
	}

	db := dbOpen()
	defer db.Close()

	note, err := db.GetNote(id)
	if err != nil {
		quitError("db move", err)
	}

	if note == nil {
		quit("db returned nil ptr :(")
	}

	edited, err := openInEditor(note.Content)
	if err != nil {
		quitError("open in editor", err)
	}

	if edited == note.Content {
		fmt.Fprintln(os.Stderr, "No changes")
		os.Exit(2)
	}

	if err = db.ReplaceContent(note.ID, edited); err != err {
		quitError("db replace", err)
	}

	fmt.Println("Note modified")
}

func checkEdit(args []string) (int, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("requires positional argument id")
	}
	return strconv.Atoi(args[0])
}
