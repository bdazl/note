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

func notePin(cmd *cobra.Command, args []string) {
	pin(args, true)
}

func noteUnpin(cmd *cobra.Command, args []string) {
	pin(args, false)
}

func pin(args []string, pinned bool) {
	ids, err := parseIds(args)
	if err != nil {
		quitError("parse ids", err)
	}

	db := dbOpen()
	defer db.Close()

	uniqueIds := removeDuplicates(ids)
	if err = db.PinNotes(uniqueIds, pinned); err != nil {
		quitError("db pin", err)
	}

	pinStr := "unpinned"
	if pinned {
		pinStr = "pinned"
	}

	count := len(uniqueIds)
	if count == 1 {
		fmt.Printf("Note %v\n", pinStr)
	} else {
		fmt.Printf("%v notes %v\n", count, pinStr)
	}
}
