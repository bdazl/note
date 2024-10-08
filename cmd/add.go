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
	"io"
	"os"
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

func noteAdd(cmd *cobra.Command, args []string) {
	note := produceNote(args)

	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	add := db.Note{
		Title:      emptyToNil(&title),
		Tags:       emptyToNil(&tags),
		Note:       note,
		IsFavorite: favorite,
	}

	id, err := db.AddNote(d, add, false)
	if err != nil {
		quitError("db add", err)
	}

	fmt.Println(id)
}

func produceNote(args []string) string {
	fileptr, err := checkAddArguments(args)
	if err != nil {
		quitError("args", err)
	}

	// Any fileptr takes precedence and we know from the arg check that len(args) == 0
	if fileptr != nil {
		data, err := io.ReadAll(fileptr)
		if err != nil {
			quitError("read file", err)
		}
		return string(data)
	}

	return noteFromArgs(args)
}

func checkAddArguments(args []string) (*os.File, error) {
	if file == "" {
		return nil, nil
	}
	// file != ""
	if len(args) > 0 {
		return nil, fmt.Errorf("you can't specify both --file and positional arguments")
	}
	if file == "-" {
		return os.Stdin, nil
	}
	return os.Open(file)
}

func noteFromArgs(args []string) string {
	if len(args) == 0 {
		// TODO: Open $EDITOR
		fmt.Fprintln(os.Stderr, "no note value provided")
		os.Exit(1)
	}

	return strings.Join(args, " ")
}

func emptyToNil(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
