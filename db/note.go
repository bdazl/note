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
	"time"
)

const (
	ColumnID         NoteColumn = "id"
	ColumnCreatedAt  NoteColumn = "created_at"
	ColumnUpdatedAt  NoteColumn = "updated_at"
	ColumnTitle      NoteColumn = "title"
	ColumnTags       NoteColumn = "tags"
	ColumnNote       NoteColumn = "note"
	ColumnIsArchived NoteColumn = "is_archived"
	ColumnIsFavorite NoteColumn = "is_favorite"
)

// Exported definition of a Note, in the DB
type Note struct {
	ID         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      *string
	Tags       *string
	Note       string
	IsArchived bool
	IsFavorite bool
}

type NoteColumn string

// Internal representation of a note
type dbNote struct {
	ID         int
	CreatedAt  string
	UpdatedAt  string
	Title      sql.NullString
	Tags       sql.NullString
	Note       string
	IsArchived bool
	IsFavorite bool
}

// Conversion helpers

func toNote(note dbNote) (*Note, error) {
	createdAt, err := parseTime(note.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := parseTime(note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &Note{
		ID:         note.ID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Title:      nullStrToPtr(note.Title),
		Tags:       nullStrToPtr(note.Tags),
		Note:       note.Note,
		IsArchived: note.IsArchived,
		IsFavorite: note.IsFavorite,
	}, nil
}

func toDbNote(note Note) dbNote {
	return dbNote{
		ID:         note.ID,
		CreatedAt:  note.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  note.UpdatedAt.Format("2006-01-02 15:04:05"),
		Title:      ptrToNullStr(note.Title),
		Tags:       ptrToNullStr(note.Tags),
		Note:       note.Note,
		IsArchived: note.IsArchived,
		IsFavorite: note.IsFavorite,
	}
}

func parseTime(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func nullStrToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func ptrToNullStr(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}
