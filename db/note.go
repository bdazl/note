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
	"strings"
	"time"
)

const (
	IDColumn         Column = "id"
	SpaceColumn      Column = "space"
	CreatedColumn    Column = "created"
	LastUpdateColumn Column = "last_updated"
	ContentColumn    Column = "content"
	PinnedColumn     Column = "pinned"
)

var (
	allNoteColumns = allNoteColumnsGen()
)

type Column string

// Exported definition of a Note, in the DB
type Note struct {
	ID         int
	Space      string
	Created    time.Time
	LastUpdate time.Time
	Content    string
	Pinned     bool
}

type Notes []Note

func (n Notes) GetIDs() []int {
	out := make([]int, len(n))
	for n, note := range n {
		out[n] = note.ID
	}
	return out
}

// Internal representation of a note
type dbNote struct {
	ID         int
	Space      string
	Created    string
	LastUpdate string
	Content    string
	Pinned     bool
}

// Helpers

func allNoteColumnsGen() string {
	// id, space, created, last_updated, content, pinned
	cols := []string{
		string(IDColumn),
		string(SpaceColumn),
		string(CreatedColumn),
		string(LastUpdateColumn),
		string(ContentColumn),
		string(PinnedColumn),
	}

	return strings.Join(cols, ", ")
}

func toNote(note dbNote) (*Note, error) {
	createdAt, err := parseTime(note.Created)
	if err != nil {
		return nil, err
	}

	updatedAt, err := parseTime(note.LastUpdate)
	if err != nil {
		return nil, err
	}

	return &Note{
		ID:         note.ID,
		Space:      note.Space,
		Created:    createdAt,
		LastUpdate: updatedAt,
		Content:    note.Content,
		Pinned:     note.Pinned,
	}, nil
}

func toDbNote(note Note) dbNote {
	return dbNote{
		ID:         note.ID,
		Space:      note.Space,
		Created:    note.Created.Format("2006-01-02 15:04:05"),
		LastUpdate: note.LastUpdate.Format("2006-01-02 15:04:05"),
		Content:    note.Content,
		Pinned:     note.Pinned,
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
