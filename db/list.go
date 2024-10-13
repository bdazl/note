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

const (
	NoLimit = 0
)

var (
	validSortColumns = map[Column]bool{
		IDColumn:         true,
		SpaceColumn:      true,
		CreatedColumn:    true,
		LastUpdateColumn: true,
		ContentColumn:    true,
		PinnedColumn:     true,
	}
)

type SortOpts struct {
	Ascending  bool
	SortColumn Column
}

func DefaultSortOpts() SortOpts {
	return SortOpts{
		Ascending:  true,
		SortColumn: IDColumn,
	}
}

func (s *SortOpts) Check() error {
	if !validSortColumns[s.SortColumn] {
		return fmt.Errorf("invalid sort column: %v", s.SortColumn)
	}
	return nil
}

func (s *SortOpts) orderStr() string {
	if s.Ascending {
		return "ASC"
	}
	return "DESC"
}

type PageOpts struct {
	Limit  int
	Offset int
}

func DefaultPageOpts() PageOpts {
	return PageOpts{}
}

func (p *PageOpts) Check() error {
	if p.Limit < 0 {
		return fmt.Errorf("limit must be positive")
	} else if p.Limit == 0 && p.Offset > 0 {
		return fmt.Errorf("if limit is zero, offset should be zero")
	}
	if p.Offset < 0 {
		return fmt.Errorf("offset must be positive")
	}
	return nil
}

func (d *DB) ListNotes(spaces []string, sortOpts *SortOpts, pageOpts *PageOpts) ([]Note, error) {
	var (
		addParams     = []any{}
		sortQueryAdd  = "ORDER BY pinned DESC" // By default we always sort pinned first
		pageQueryAdd  = ""
		spaceQueryAdd = ""
		limit         = 0
	)
	// Check input
	if sortOpts != nil {
		if err := sortOpts.Check(); err != nil {
			return nil, err
		}
		order := sortOpts.orderStr()
		sortColumn := string(sortOpts.SortColumn)
		sortQueryAdd = fmt.Sprintf("ORDER BY pinned DESC, %v %v", sortColumn, order)
	}
	if pageOpts != nil {
		if err := pageOpts.Check(); err != nil {
			return nil, err
		}
		if pageOpts.Limit > 0 {
			pageQueryAdd = "LIMIT ? OFFSET ?"
			addParams = append(addParams, pageOpts.Limit)
			addParams = append(addParams, pageOpts.Offset)
		}
		limit = pageOpts.Limit
	}
	if len(spaces) > 0 {
		spaceQueryAdd = spacesWhere(len(spaces))
		for _, space := range spaces {
			addParams = append(addParams, space)
		}
	}

	// Prepare the SQL query
	query := fmt.Sprintf(
		"SELECT %v FROM notes %v %v %v",
		allNoteColumns, spaceQueryAdd, sortQueryAdd, pageQueryAdd,
	)

	// Execute the query
	rows, err := d.db.Query(query, addParams...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Parse the results
	notes := make([]Note, 0, limit)
	for rows.Next() {
		note, err := scanNote(rows)
		if err != err {
			return nil, err
		}
		notes = append(notes, *note)
	}

	return notes, nil
}

func (d *DB) ListSpaces(sortOpts *SortOpts) ([]string, error) {
	orderBy := ""
	if sortOpts != nil {
		orderBy = fmt.Sprintf("ORDER BY %v %v", sortOpts.SortColumn, sortOpts.orderStr())
	}

	query := fmt.Sprintf("SELECT DISTINCT space FROM notes %v", orderBy)
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	spaces := make([]string, 0)
	for rows.Next() {
		var space string
		err := rows.Scan(&space)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		spaces = append(spaces, space)
	}
	return spaces, nil
}

func spacesWhere(count int) string {
	if count < 1 {
		return ""
	}
	return fmt.Sprintf("WHERE %v", equalOrChain("space", count))
}

func equalOrChain(lhs string, count int) string {
	var (
		bld      = strings.Builder{}
		orStmt   = " OR "
		mainStmt = fmt.Sprintf("%v = ?", lhs)
	)
	for i := 0; i < count; i++ {
		if i > 0 {
			bld.WriteString(orStmt)
		}
		bld.WriteString(mainStmt)
	}
	return bld.String()
}
