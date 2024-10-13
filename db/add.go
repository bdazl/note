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
)

// Add a note to the database.
// If full is true, then all values (except ID) are taken from the input,
// otherwise timestamps and other default values are set automatically.
func (d *DB) AddNote(note Note, full bool) (int64, error) {
	const (
		smallQuery = "INSERT INTO notes (space, content, pinned) VALUES (?, ?, ?);"
		fullQuery  = `INSERT INTO notes (space, created, last_updated, content, pinned)
		VALUES (?, ?, ?, ?, ?);`
	)
	var (
		dbN    = toDbNote(note)
		query  string
		params []any
	)

	if full {
		query = fullQuery
		params = []any{
			dbN.Space,
			dbN.Created,
			dbN.LastUpdate,
			dbN.Content,
			dbN.Pinned,
		}
	} else {
		query = smallQuery
		params = []any{dbN.Space, dbN.Content, dbN.Pinned}
	}

	result, err := d.db.Exec(query, params...)
	if err != nil {
		return 0, fmt.Errorf("insert error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("last insert id error: %w", err)
	}

	return id, nil
}
