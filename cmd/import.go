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
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func noteImport(cmd *cobra.Command, args []string) {
	paths, err := uniquePaths(args)
	if err != nil {
		quitError("args", err)
	}

	preferredFmt, err := cmdArgFormat()
	if err != nil {
		quitError("args", err)
	}

	allNotes := make([]FileNote, 0)
	for _, path := range paths {
		fileFmt := filenameFormat(path)
		if fileFmt == UnknownFormat {
			// We know that if we get here the preferredFmt should not be unknown
			// This is checked during args validation
			fileFmt = preferredFmt
		}

		reader, err := openFile(path)
		if err != nil {
			quitError("open file", err)
		}
		defer reader.Close()

		var notes []FileNote
		switch fileFmt {
		case JSONFormat:
			notes, err = decodeJSON(reader)
			if err != nil {
				quitError("decode JSON", err)
			}
		case YAMLFormat:
			notes, err = decodeYAML(reader)
			if err != nil {
				quitError("decode YAML", err)
			}
		default:
			quit("unknown format")
		}

		allNotes = append(allNotes, notes...)
	}

	dbNotes := fileNotesToDB(allNotes)

	db := dbOpen()
	defer db.Close()

	ids := make([]int, 0, len(dbNotes))
	for _, note := range dbNotes {
		id, err := db.AddNote(note, true)
		if err != nil {
			err = fmt.Errorf(
				"generated ids: %v, but only %v of %v successful. db error: %w",
				ids, len(ids), len(dbNotes), err)
			quitError("db add", err)
		}
		ids = append(ids, int(id))
	}

	if listArg {
		for _, id := range ids {
			fmt.Println(id)
		}
	} else {
		idStrs := manyIntToString(ids)
		joined := strings.Join(idStrs, ", ")
		fmt.Printf("Notes created: %v\n", joined)
	}
}

func fileNotesToDB(notes []FileNote) db.Notes {
	out := make(db.Notes, len(notes))
	for n, note := range notes {
		out[n] = db.Note{
			Pinned:      note.Pinned,
			Space:       note.Space,
			Content:     note.Content,
			Created:     note.Created,
			LastUpdated: note.LastUpdated,
		}
	}
	return out
}

func decodeJSON(reader io.Reader) ([]FileNote, error) {
	decoder := json.NewDecoder(reader)
	out := make([]FileNote, 0)
	for {
		var notes []FileNote
		if err := decoder.Decode(&notes); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		out = append(out, notes...)
	}
	return out, nil
}

func decodeYAML(reader io.Reader) ([]FileNote, error) {
	var notes []FileNote

	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(&notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func uniquePaths(args []string) ([]string, error) {
	unique := removeDuplicates(args)
	for _, path := range unique {
		if path != StdoutPath {
			if !exists(path) {
				return nil, fmt.Errorf("path does not exist: %v", path)
			}
		} else {
			preferredFmt, err := cmdArgFormat()
			if err != nil {
				return nil, fmt.Errorf("cmd fmt error: %w", err)
			}
			if preferredFmt == UnknownFormat {
				return nil, fmt.Errorf("if stdin is used a preferred file format must be chosen")
			}
		}
	}
	return unique, nil
}

func manyIntToString(ints []int) []string {
	out := make([]string, len(ints))
	for n, val := range ints {
		out[n] = strconv.Itoa(val)
	}
	return out
}
