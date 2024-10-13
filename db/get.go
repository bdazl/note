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
