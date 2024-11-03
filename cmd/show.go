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

func noteShow(cmd *cobra.Command, args []string) {
	ids, err := parseIds(args)
	if err != nil {
		quitError("parse ids", err)
	}

	style, color, err := styleColorOpts()
	if err != nil {
		quitError("args", err)
	}

	db := dbOpen()
	defer db.Close()

	uniqueIds := removeDuplicates(ids)
	notes, err := db.GetNotes(uniqueIds)
	if err != nil {
		quitError("db get", err)
	}

	// Minimal map of notes
	notesMap := notes.AsMap()

	// Expand user input order here. If she wanted the same id multiple times, she shall have it.
	notesOrdered := expandNoteIds(notesMap, ids)

	// This command has a special case for when the user requested to get one note
	// For this case, default to only printing the content of the note.
	if len(notesOrdered) == 1 && !alwaysStyleArg {
		printNotesMinimal(notesOrdered)
	} else {
		pprintNotes(notesOrdered, style, color)
	}
}

func expandNoteIds(noteMap db.NoteMap, ids []int) db.Notes {
	out := make(db.Notes, len(ids))
	for i, id := range ids {
		var ok bool
		out[i], ok = noteMap[id]
		if !ok {
			quit("internal error: mapping non-existing note id")
		}
	}
	return out
}

func parseIds(args []string) ([]int, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("require at least one id")
	}

	ids := make([]int, len(args))
	for i, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("parse error for %v: %w", arg, err)
		}
		ids[i] = id
	}

	return ids, nil
}

// removeDuplicates but preserve order
func removeDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, value := range slice {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}
