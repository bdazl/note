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
		Run:   noteList,
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
		Use:   "add [notes...]",
		Short: "Add new note",
		Run:   noteAdd,
	}
	rmCmd = &cobra.Command{
		Use:     "remove [id...]",
		Aliases: []string{"rm"},
		Short:   "Remove note(s) with id(s)",
		Run:     noteRemove,
	}
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Prints some or all of your notes",
		Run:     noteList,
	}

	// Global arguments
	configPath         string
	storagePathCmdLine string

	// Init argument
	force bool

	// Add arguments
	title    string
	tags     string
	file     string
	favorite bool

	// List arguments
	sortBy     string
	descending bool
	limit      int
	offset     int
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		quitError("root exec", err)
	}
}

func init() {
	dfltConfig, err := defaultConfigPath()
	if err != nil {
		quitError("default config path", err)
	}
	dfltStore, err := defaultStoragePath()
	if err != nil {
		quitError("default storage path", err)
	}

	globalFlags := rootCmd.PersistentFlags()
	globalFlags.StringVarP(&configPath, "config", "c", dfltConfig, "config file")
	globalFlags.StringVar(&storagePathCmdLine, "db", dfltStore, "database store containing your notes")

	// db can exist in config file
	viper.BindPFlag("db", globalFlags.Lookup("db"))

	initFlags := initCmd.Flags()
	initFlags.BoolVar(&force, "force", false, "determines if existing files will be overwritten")

	addFlags := addCmd.Flags()
	addFlags.StringVarP(&title, "name", "n", "", "title of note")
	addFlags.StringVarP(&tags, "tags", "t", "", "tags of note as a comma separated string")
	addFlags.StringVarP(&file, "file", "f", "", "the note is read from file")
	addFlags.BoolVar(&favorite, "fav", false, "mark note as favorite")

	sortKeys := getSortKeys()
	sortUsage := fmt.Sprintf("column to sort notes by (%v)", sortKeys)
	listFlags := listCmd.Flags()
	listFlags.StringVarP(&sortBy, "sort", "s", "id", sortUsage)
	listFlags.IntVarP(&limit, "limit", "l", 0, "limit amount of notes shown, 0 means no limit")
	listFlags.IntVarP(&offset, "offset", "o", 0, "begin list notes at some offset")
	listFlags.BoolVarP(&descending, "descend", "r", false, "descending order")

	// TODO: It makes not sense to include all of lists arguments here
	// rootCmd.Flags().AddFlagSet(listFlags)

	rootCmd.AddCommand(initCmd, addCmd, rmCmd, listCmd)
}

func noteInit(cmd *cobra.Command, args []string) {
	forceInform := false
	dbF, err := filepath.Abs(storagePathCmdLine) // When doing init we explicitly want the command line option
	if err != nil {
		quitError("db path", err)
	}

	// We force viper to set value to the (abs path of the) command line option here
	// This is because init may be re-ran after we have a valid config setup, and we don't want to then source
	// the option from the config file.
	viper.Set("db", dbF)

	mkdir(filepath.Dir(configPath))
	mkdir(filepath.Dir(dbF))

	if !force && exists(configPath) {
		fmt.Fprintln(os.Stderr, "Config file already exists")
		forceInform = true
	} else {
		fmt.Printf("Writing config file: %v\n", configPath)
		err := viper.WriteConfig()
		if err != nil {
			quitError("writing config", err)
		}
	}

	if !force && exists(dbF) {
		fmt.Fprintln(os.Stderr, "Storage file already exists")
		forceInform = true
	} else {
		fmt.Printf("Create initial db: %v\n", dbF)
		if _, err := db.CreateDb(dbF); err != nil {
			quitError("creating db", err)
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
		quitError("mkdir", err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func quitError(loc string, err error) {
	msg := fmt.Sprintf("error %v: %v", loc, err)
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
