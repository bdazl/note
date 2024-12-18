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
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func noteAdd(cmd *cobra.Command, args []string) {
	content := produceNote(args)

	space := viper.GetString(ViperSpace)
	if err := checkSpaceArgument(space); err != nil {
		quitError("arg", err)
	}

	add := db.Note{
		Space:   viper.GetString(ViperSpace),
		Content: content,
		Pinned:  pinnedArg,
	}

	d := dbOpen()
	defer d.Close()

	id, err := d.AddNote(add, false)
	if err != nil {
		quitError("db add", err)
	}

	fmt.Printf("Created note: %v\n", id)
}

func produceNote(args []string) string {
	reader, err := checkAddArguments(args)
	if err != nil {
		quitError("args", err)
	}

	// Any fileptr takes precedence and we know from the arg check that len(args) == 0
	if reader != nil {
		defer reader.Close()
		return readAll(reader)
	}

	// If there are arguments, interpret them as strings with spaces in between
	if len(args) != 0 {
		return strings.Join(args, " ")
	}

	// Special case, where no arguments means open an editor to create the note
	note, err := openInEditor("")
	if err != nil {
		quitError("open editor", err)
	}
	return note
}

func checkSpaceArgument(space string) error {
	if strings.Contains(space, ",") {
		return fmt.Errorf("space cannot contain the following character ','")
	}
	return nil
}

func checkAddArguments(args []string) (io.ReadCloser, error) {
	if fileArg == "" {
		return nil, nil
	}
	if len(args) > 0 {
		return nil, fmt.Errorf("you can't specify both --file and positional arguments")
	}
	return openFile(fileArg)
}
