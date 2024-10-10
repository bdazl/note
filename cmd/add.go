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
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func noteAdd(cmd *cobra.Command, args []string) {
	content := produceNote(args)

	space := viper.GetString(ViperAddSpace)
	if err := checkSpaceArgument(space); err != nil {
		quitError("arg", err)
	}

	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	add := db.Note{
		Space:    viper.GetString(ViperAddSpace),
		Content:  content,
		IsPinned: pinnedArg,
	}

	id, err := d.AddNote(add, false)
	if err != nil {
		quitError("db add", err)
	}

	fmt.Println(id)
}

func produceNote(args []string) string {
	fileptr, err := checkAddArguments(args)
	if err != nil {
		quitError("args", err)
	}

	// Any fileptr takes precedence and we know from the arg check that len(args) == 0
	if fileptr != nil {
		data, err := io.ReadAll(fileptr)
		if err != nil {
			quitError("read file", err)
		}
		return string(data)
	}

	if len(args) == 0 {
		note, err := openInEditor()
		if err != nil {
			quitError("open editor", err)
		}
		return note
	}

	return strings.Join(args, " ")
}

func openInEditor() (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "note.*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the file is removed

	// Close the file so the editor can open it
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	// Get the editor from configuration or environment
	editor := viper.GetString("editor")

	// Open the file with the editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to open editor: %w", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read file after close: %w", err)
	}

	if len(content) == 0 {
		return "", fmt.Errorf("nothing to add")
	}

	return string(content), nil
}

func checkSpaceArgument(space string) error {
	if strings.Contains(space, ",") {
		return fmt.Errorf("space cannot contain the following character ','")
	}
	return nil
}

func checkAddArguments(args []string) (*os.File, error) {
	if fileArg == "" {
		return nil, nil
	}
	// file != ""
	if len(args) > 0 {
		return nil, fmt.Errorf("you can't specify both --file and positional arguments")
	}
	if fileArg == "-" {
		return os.Stdin, nil
	}
	return os.Open(fileArg)
}
