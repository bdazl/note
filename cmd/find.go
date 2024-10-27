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
package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
)

type finder struct {
	*regexp.Regexp
	Expression  string
	Insensitive bool
}

func finderFromArgs(args []string) (*finder, error) {
	var (
		expr  = strings.Join(args, " ")
		err   error
		regex *regexp.Regexp
	)

	if regexpArg && !posixArg {
		regex, err = regexp.Compile(expr)
		if err != nil {
			return nil, fmt.Errorf("compile: %w", err)
		}
	} else if posixArg {
		regex, err = regexp.CompilePOSIX(expr)
		if err != nil {
			return nil, fmt.Errorf("compile: %w", err)
		}
	}

	return &finder{
		Regexp:      regex,
		Expression:  expr,
		Insensitive: insensitiveArg,
	}, nil
}

func (f *finder) Match(str string) bool {
	if f.Regexp != nil {
		return f.Regexp.MatchString(str)
	} else if !f.Insensitive {
		return strings.Contains(str, f.Expression)
	} else {
		// Make this more efficient
		lower := strings.ToLower(str)
		return strings.Contains(lower, f.Expression)
	}
}

func noteFind(cmd *cobra.Command, args []string) {
	finder, err := finderFromArgs(args)
	if err != nil {
		quitError("pattern", err)
	}

	style, color, err := styleColorOpts()
	if err != nil {
		quitError("args", err)
	}

	d := dbOpen()
	defer d.Close()

	notes := make(db.Notes, 0)
	for iter := range d.IterateNotes(nil, allArg || trashArg, nil) {
		if iter.Err != nil {
			quitError("db iterate", iter.Err)
		}

		if !trashArg && iter.Note.Space == TrashSpace {
			continue
		}

		if finder.Match(iter.Note.Content) {
			notes = append(notes, iter.Note)
		}
	}

	if idArg {
		ids := notes.GetIDs()
		printIds(ids, listArg)
	} else {
		pprintNotes(notes, style, color)
	}
}
