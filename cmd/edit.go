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
	"os"
	"os/exec"
	"strconv"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func noteEdit(cmd *cobra.Command, args []string) {
	id, err := checkEdit(args)
	if err != nil {
		quitError("args", err)
	}

	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}

	note, err := d.GetNote(id)
	if err != nil {
		quitError("db move", err)
	}

	if note == nil {
		quit("db returned nil ptr :(")
	}

	edited, err := openInEditor(note.Content)
	if err != nil {
		quitError("open in editor", err)
	}

	if edited == note.Content {
		fmt.Fprintln(os.Stderr, "No changes")
		os.Exit(2)
	}

	if err = d.ReplaceContent(note.ID, edited); err != err {
		quitError("db replace", err)
	}

	fmt.Println("Note modified")
}

func openInEditor(initText string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "note.*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the file is removed

	// Write initial text to temp file, if any
	if initText != "" {
		_, err := tmpFile.WriteString(initText)
		if err != nil {
			return "", fmt.Errorf("failed to write to temp file: %w", err)
		}
	}

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

	// Read what the user wrote
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read file after close: %w", err)
	}

	// If nothing exists in the file we always consider it an error
	if len(content) == 0 {
		return "", fmt.Errorf("empty file")
	}

	return string(content), nil
}

func checkEdit(args []string) (int, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("requires positional argument id")
	}
	return strconv.Atoi(args[0])
}
