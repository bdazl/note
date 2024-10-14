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
	"strings"
)

func (d *DB) ReplaceContent(id int, content string) error {
	return d.updateRow("UPDATE notes SET content = ? WHERE id = ?", content, id)
}

func (d *DB) MoveNote(id int, toSpace string) error {
	return d.updateRow("UPDATE notes SET space = ? WHERE id = ?", toSpace, id)
}

func (d *DB) MoveNotes(ids []int, toSpace string) error {
	count := len(ids)
	if count < 1 {
		return fmt.Errorf("require at least one id")
	}

	// ID slots
	manyQuestions := repeatString("?", count)
	bracketQ := strings.Join(manyQuestions, ", ")

	// Construct query
	query := fmt.Sprintf(
		"UPDATE notes SET space = ? WHERE id IN (%v)",
		bracketQ,
	)

	// Execute
	execParams := sliceToAny(ids)
	execParams = prepend(execParams, any(toSpace))
	result, err := d.db.Exec(query, execParams...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows: %w", err)
	}

	// Validate
	if rows != int64(count) {
		return fmt.Errorf("only %v out of %v was moved successfully", rows, count)
	}

	return nil
}

func (d *DB) PinNotes(ids []int, pinned bool) error {
	count := len(ids)
	if count < 1 {
		return fmt.Errorf("must provide ids")
	}

	pinVal := "0"
	if pinned {
		pinVal = "1"
	}

	// ID slots
	manyQuestions := repeatString("?", count)
	bracketQ := strings.Join(manyQuestions, ", ")

	// Construct query
	query := fmt.Sprintf(
		"UPDATE notes SET pinned = %v WHERE id IN (%v)",
		pinVal, bracketQ,
	)

	// Execute
	result, err := d.db.Exec(query, sliceToAny(ids)...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows: %w", err)
	}

	// Validate
	if rows != int64(count) {
		return fmt.Errorf("only %v out of %v was pinned successfully", rows, count)
	}

	return nil
}

func (d *DB) updateRow(query string, args ...any) error {
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("nothing was modified")
	}

	return nil
}

func prepend[T any](slice []T, element T) []T {
	return append([]T{element}, slice...)
}
