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
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

func noteAdd(cmd *cobra.Command, args []string) {
	d, err := db.Open(storagePath)
	if err != nil {
		panic(err)
	}

	add := db.Note{
		Title:      emptyToNil(&title),
		Tags:       emptyToNil(&tags),
		Note:       argsToNote(args),
		IsFavorite: favorite,
	}

	id, err := db.AddNote(d, add, false)
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
}

func argsToNote(args []string) string {
	if len(args) == 0 {
		// TODO: Open $EDITOR
		fmt.Fprintln(os.Stderr, "no note value provided")
		os.Exit(1)
	}

	sep := " "
	if newline {
		sep = "\n"
	}
	return strings.Join(args, sep)
}

func emptyToNil(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
