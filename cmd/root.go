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
	"path/filepath"

	"github.com/bdazl/note/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	currentCmd *cobra.Command

	rootCmd = &cobra.Command{
		Use:   "note",
		Short: "No fuzz terminal note taking",
		Run:   noteRoot,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			currentCmd = cmd
			initConfig()
		},
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize note configuration",
		Run:   noteInit,
	}
	addCmd = &cobra.Command{
		Use:     "add note [notes...]",
		Aliases: []string{"a"},
		Short:   "Add a new note",
		Run:     noteAdd,
	}
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Prints some or all of your notes",
		Run:     noteList,
	}
)

// Command line argument values
var (
	// Global arguments
	configPath  string
	storagePath string

	// Init argument
	force bool

	// Add arguments
	title     string
	tags      string
	separator string
	favorite  bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err.Error())
		os.Exit(1)
	}
}

func init() {
	dfltConfig, err := defaultConfigPath()
	if err != nil {
		panic(err)
	}
	dfltStore, err := defaultStoragePath()
	if err != nil {
		panic(err)
	}

	globalFlags := rootCmd.PersistentFlags()
	globalFlags.StringVarP(&configPath, "config", "c", dfltConfig, "config file")
	globalFlags.StringVarP(&storagePath, "store", "d", dfltStore, "database store containing your notes")

	initFlags := initCmd.Flags()
	initFlags.BoolVar(&force, "force", false, "determines if existing files will be overwritten")

	addFlags := addCmd.Flags()
	addFlags.StringVarP(&title, "name", "n", "", "title of note (optional)")
	addFlags.StringVarP(&tags, "tags", "t", "", "tags of note as a comma separated string (optional)")
	addFlags.StringVarP(&separator, "separator", "s", " ", "concatenate positional arguments with this separator")
	addFlags.BoolVarP(&favorite, "favorite", "f", false, "mark note as favorite")

	rootCmd.AddCommand(initCmd, addCmd, listCmd)
}

func noteInit(cmd *cobra.Command, args []string) {
	forceInform := false

	mkdir(filepath.Dir(configPath))
	mkdir(filepath.Dir(storagePath))

	if !force && exists(configPath) {
		fmt.Fprintln(os.Stderr, "Config file already exists")
		forceInform = true
	} else {
		fmt.Printf("Writing config file: %v\n", configPath)
		err := viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	}

	if !force && exists(storagePath) {
		fmt.Fprintln(os.Stderr, "Storage file already exists")
		forceInform = true
	} else {
		fmt.Printf("Create initial db: %v\n", storagePath)
		if err := db.CreateDb(storagePath); err != nil {
			panic(err)
		}
	}

	if forceInform {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "Some file(s) where not initialized")
		fmt.Fprintln(os.Stderr, "If you want to force re-create them, consider using the --force flag")
	}
}

func mkdir(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
