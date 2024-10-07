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
	"database/sql"
	"fmt"
)

const (
	NoLimit = 0
)

func ListNotes(db *sql.DB, sortBy NoteColumn, ascending bool, limit int, offset int) ([]Note, error) {
	err := validSortColumn(sortBy)
	if err != nil {
		return nil, err
	}
	if limit < 0 {
		return nil, fmt.Errorf("limit must be positive")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset must be positive")
	}

	// Determine sort order
	order := "ASC"
	if !ascending {
		order = "DESC"
	}

	// Prepare the SQL query
	loff := ""
	addParams := []any{}
	if limit > 0 {
		loff = "LIMIT ? OFFSET ?"
		addParams = append(addParams, limit)
		addParams = append(addParams, offset)
	}
	query := fmt.Sprintf(`
        SELECT id, created_at, updated_at, title, tags, note, is_archived, is_favorite
        FROM notes
        ORDER BY %v %v
        %v`, string(sortBy), order, loff,
	)

	// Execute the query
	rows, err := db.Query(query, addParams...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Parse the results
	notes := make([]Note, 0, limit)
	var dbN dbNote
	for rows.Next() {
		err := rows.Scan(
			&dbN.ID,
			&dbN.CreatedAt,
			&dbN.UpdatedAt,
			&dbN.Title,
			&dbN.Tags,
			&dbN.Note,
			&dbN.IsArchived,
			&dbN.IsFavorite,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		nOut, err := toNote(dbN)
		if err != nil {
			return nil, fmt.Errorf("row conversion error: %w", err)
		}

		notes = append(notes, *nOut)
	}

	return notes, nil
}

func validSortColumn(col NoteColumn) error {
	validSortColumns := map[NoteColumn]bool{
		ColumnID:         true,
		ColumnCreatedAt:  true,
		ColumnUpdatedAt:  true,
		ColumnTitle:      true,
		ColumnTags:       true,
		ColumnNote:       true,
		ColumnIsArchived: true,
		ColumnIsFavorite: true,
	}

	if !validSortColumns[col] {
		return fmt.Errorf("invalid sort column: %v", col)
	}
	return nil
}
