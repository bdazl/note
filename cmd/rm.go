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

func noteRemove(cmd *cobra.Command, args []string) {
	ids, err := checkRemoveArguments(args)
	if err != nil {
		quitError("args", err)
	}

	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	for _, id := range ids {
		if err := d.RemoveNote(id); err != nil {
			quitError("db remove", err)
		}
	}
}

func checkRemoveArguments(args []string) ([]int, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("you must provide at least one id")
	}

	out := make([]int, len(args))
	for i, ids := range args {
		id, err := strconv.Atoi(ids)
		if err != nil {
			return nil, err
		}
		out[i] = id
	}
	return out, nil
}
