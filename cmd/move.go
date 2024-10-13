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
	"strconv"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

func noteMove(cmd *cobra.Command, args []string) {
	id, toSpace, err := checkMove(args)
	if err != nil {
		quitError("args", err)
	}

	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	if err = d.MoveNote(id, toSpace); err != nil {
		quitError("db move", err)
	}

	fmt.Println("Note modified")
}

func checkMove(args []string) (int, string, error) {
	if len(args) != 2 {
		return 0, "", fmt.Errorf("requires positional arguments id and toSpace")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, "", fmt.Errorf("id parse error: %w", err)
	}
	return id, args[1], nil
}
