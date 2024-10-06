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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	currentCmd *cobra.Command

	rootCmd = &cobra.Command{
		Use:   "note",
		Short: "No fuzz terminal note taking",
		Run:   list, // list notes per default
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			currentCmd = cmd
			initConfig()
		},
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize note configuration",
		Run:   noteInit, // list notes per default
	}
	addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new note",
		Run:     list, // list notes per default
	}
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Prints some or all of your notes",
		Run:     list, // list notes per default
	}
)

// Command line argument values
var (
	configPath  string
	storagePath string
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

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&configPath, "config", dfltConfig, "config file")
	pflags.StringVar(&storagePath, "store", dfltStore, "database store containing your notes")

	rootCmd.AddCommand(initCmd, addCmd, listCmd)
}

func noteInit(cmd *cobra.Command, args []string) {
	fmt.Printf("Writing config file: %v\n", configPath)

	mkdir(filepath.Dir(configPath))
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create initial db: %v\n", storagePath)

	mkdir(filepath.Dir(storagePath))
	// TODO: create storage
}

func mkdir(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
