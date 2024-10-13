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

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Scanner interface {
	Scan(dest ...any) error
}

func (d *DB) GetNote(id int) (*Note, error) {
	query := fmt.Sprintf("SELECT %v FROM notes WHERE id = ?", allNoteColumns)
	row := d.db.QueryRow(query, id)

	note, err := scanNote(row)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (d *DB) GetNotes(ids []int) ([]Note, error) {
	count := len(ids)
	if count < 1 {
		return nil, fmt.Errorf("require at least one id")
	}

	// ID slots
	manyQuestions := repeatString("?", count)
	bracketQ := strings.Join(manyQuestions, ", ")

	// The IDs must be returned in the same order
	whenChain := whenThenChain(ids)
	orderBy := fmt.Sprintf("ORDER BY CASE id %v END", whenChain)

	// Construct query from above
	query := fmt.Sprintf(
		"SELECT %v FROM notes WHERE id IN (%v) %v",
		allNoteColumns,
		bracketQ,
		orderBy,
	)

	// Query with ids as []any
	idsAsAny := sliceToAny(ids)
	rows, err := d.db.Query(query, idsAsAny...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Parse results
	notes := make([]Note, 0)
	for rows.Next() {
		note, err := scanNote(rows)
		if err != err {
			return nil, err
		}
		notes = append(notes, *note)
	}

	// If we did not get all id's, figure out which ones where not found
	if len(notes) != count {
		outIds := idsFromNotes(notes)
		diff := difference(ids, outIds)
		idStrs := manyIntToString(diff)
		joined := strings.Join(idStrs, ", ")
		return nil, fmt.Errorf("the following ids did not exist: %v", joined)
	}

	return notes, nil
}

func scanNote(scanner Scanner) (*Note, error) {
	var dbN dbNote
	err := scanner.Scan(
		&dbN.ID,
		&dbN.Space,
		&dbN.Created,
		&dbN.LastUpdate,
		&dbN.Content,
		&dbN.Pinned,
	)

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	out, err := toNote(dbN)
	if err != nil {
		return nil, fmt.Errorf("conversion error: %w", err)
	}

	return out, nil
}

func repeatString(str string, count int) []string {
	out := make([]string, count)
	for i := 0; i < count; i++ {
		out[i] = str
	}
	return out
}

func whenThenChain(ids []int) string {
	bld := strings.Builder{}

	for n, id := range ids {
		if n > 0 {
			bld.WriteRune(' ')
		}
		bld.WriteString(fmt.Sprintf("WHEN %v THEN %v", id, n))
	}

	return bld.String()
}

func manyIntToString(ints []int) []string {
	out := make([]string, len(ints))
	for n, val := range ints {
		out[n] = strconv.Itoa(val)
	}
	return out
}

func difference(full, subset []int) []int {
	out := make([]int, 0, len(full))
	for _, lhs := range full {
		if !slices.Contains(subset, lhs) {
			out = append(out, lhs)
		}
	}
	return out
}

func idsFromNotes(notes []Note) []int {
	out := make([]int, len(notes))
	for n, note := range notes {
		out[n] = note.ID
	}
	return out
}
