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

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTableSql = `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		title TEXT DEFAULT NULL,
		tags TEXT DEFAULT NULL,
		is_archived BOOLEAN DEFAULT 0,
		is_favorite BOOLEAN DEFAULT 0,
		note TEXT NOT NULL
	);
`

	createTriggerSql = `CREATE TRIGGER IF NOT EXISTS notes_auto_updated_at
		AFTER UPDATE ON notes
		FOR EACH ROW
		BEGIN
			UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
		END;
`
)

func CreateDb(path string) error {
	db, err := Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(createTableSql)
	if err != nil {
		return fmt.Errorf("create notes table: %w", err)
	}

	_, err = db.Exec(createTriggerSql)
	if err != nil {
		return fmt.Errorf("create trigger: %w", err)
	}
	return nil
}
