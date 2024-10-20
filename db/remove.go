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

func (d *DB) PermanentRemoveNotes(ids []int) error {
	count := len(ids)
	if count < 1 {
		return fmt.Errorf("must provide ids")
	}

	// ID slots
	manyQuestions := repeatString("?", len(ids))
	bracketQ := strings.Join(manyQuestions, ", ")

	idsWhere := fmt.Sprintf("id IN (%v)", bracketQ)

	// Combine query
	query := fmt.Sprintf("DELETE FROM notes WHERE %v", idsWhere)

	// Execute query
	result, err := d.db.Exec(query, sliceToAny(ids)...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows: %w", err)
	}

	if rows != int64(count) {
		return fmt.Errorf("only %v out of %v was deleted successfully", rows, count)
	}

	return nil
}

func sliceToAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}
