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
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	UnknownFormat FileFormat = iota
	JSONFormat
	YAMLFormat
)

type FileFormat int

type FileNote struct {
	ID          int       `json:"id" yaml:"id"`
	Pinned      bool      `json:"pinned" yaml:"pinned"`
	Space       string    `json:"space" yaml:"space"`
	Content     string    `json:"content" yaml:"content"`
	Created     time.Time `json:"created" yaml:"created"`
	LastUpdated time.Time `json:"last_updated" yaml:"last_updated"`
}

func noteExport(cmd *cobra.Command, args []string) {
	argFmt, err := cmdArgFormat()
	if err != nil {
		quitError("cmd arg fmt", err)
	}

	path, fileFmt, err := exportFilePathAndFormat(args)
	if err != nil {
		quitError("path arg", err)
	}

	format, err := combinedFormat(argFmt, fileFmt)
	if err != nil {
		quitError("file format", err)
	}

	notes, err := selectNotes(spacesArg)
	if err != nil {
		quitError("collect notes", err)
	}

	fileNotes := convFileNotes(notes)

	writer, err := createFileOrStdout(path)
	if err != nil {
		quitError("open writer", err)
	}
	defer writer.Close()

	if format == JSONFormat {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(jsonPrefixArg, jsonIndentArg)
		if err = encoder.Encode(fileNotes); err != nil {
			quitError("encode", err)
		}
	} else if format == YAMLFormat {
		encoder := yaml.NewEncoder(writer)
		encoder.SetIndent(yamlSpacesArg)
		if err = encoder.Encode(fileNotes); err != nil {
			quitError("encode", err)
		}
	}
}

func convFileNotes(notes db.Notes) []FileNote {
	converted := make([]FileNote, len(notes))
	for i, note := range notes {
		converted[i] = FileNote{
			ID:          note.ID,
			Pinned:      note.Pinned,
			Space:       note.Space,
			Content:     note.Content,
			Created:     note.Created,
			LastUpdated: note.LastUpdated,
		}
	}
	return converted
}

func exportFilePathAndFormat(args []string) (string, FileFormat, error) {
	if len(args) == 0 {
		return StdoutPath, UnknownFormat, nil
	} else if len(args) != 1 {
		return "", UnknownFormat, fmt.Errorf("must only contain one positional argument")
	}

	path := args[0]
	if forceArg && exists(path) {
		return "", UnknownFormat, fmt.Errorf("file already exist")
	}
	return path, filenameFormat(path), nil
}

func combinedFormat(argFmt, fileFmt FileFormat) (FileFormat, error) {
	// The final format is per default the argument format
	if argFmt != UnknownFormat {
		return argFmt, nil
	}

	// If this has not been input, then we check if the file extension offers the answer
	if fileFmt != UnknownFormat {
		return fileFmt, nil
	}

	return UnknownFormat, fmt.Errorf("could not determine output format")
}

// cmdArgFormat specifies the file format set by the command line option
func cmdArgFormat() (FileFormat, error) {
	if jsonArg {
		if yamlArg {
			return UnknownFormat, fmt.Errorf("you can only pick one export format")
		}
		return JSONFormat, nil
	} else if yamlArg {
		return YAMLFormat, nil
	}
	return UnknownFormat, nil
}

func filenameFormat(path string) FileFormat {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return JSONFormat
	case ".yml":
		return YAMLFormat
	case ".yaml":
		return YAMLFormat
	}
	return UnknownFormat
}
