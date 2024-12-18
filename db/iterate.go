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
package db

import "fmt"

type NoteIterator struct {
	Note
	Err error
}

func (d *DB) IterateNotes(spaces []string, all bool, sortOpts *SortOpts) <-chan NoteIterator {
	ch := make(chan NoteIterator)

	go func() {
		defer close(ch)

		page := 0
		pageOpts := &PageOpts{
			Limit: 10,
		}

		for {
			notes, err := d.SelectNotes(spaces, all, sortOpts, pageOpts)
			if err != nil {
				fmt.Printf("Error: %v", err.Error())
				ch <- NoteIterator{Err: err}
				return
			}

			if len(notes) == 0 {
				return
			}

			for _, note := range notes {
				ch <- NoteIterator{Note: note}
			}

			page += 1
			pageOpts.Offset = pageOpts.Limit * page
		}
	}()

	return ch
}
