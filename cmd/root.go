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

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	removeCmd = &cobra.Command{
		Use:     "remove id [id...]",
		Aliases: []string{"rm", "del"},
		Short:   "Remove note(s) with id(s)",
		Run:     noteRemove,
	}
	getCmd = &cobra.Command{
		Use:   "get id [id...]",
		Short: "Get specific note(s)",
		Args:  cobra.MinimumNArgs(1),
		Run:   noteGet,
	}
	listCmd = &cobra.Command{
		Use:     "list [space...]",
		Aliases: []string{"ls"},
		Short:   "Lists notes from one or more spaces",
		Run:     noteList,
	}
	editCmd = &cobra.Command{
		Use:   "edit id",
		Short: "Edit content of note",
		Args:  cobra.MinimumNArgs(1),
		Run:   noteEdit,
	}
	pinCmd = &cobra.Command{
		Use:   "pin id [id...]",
		Short: "Pin note(s) to top",
		Args:  cobra.MinimumNArgs(1),
		Run:   notePin,
	}
	unpinCmd = &cobra.Command{
		Use:   "unpin id [id...]",
		Short: "Unpin note(s) from top",
		Args:  cobra.MinimumNArgs(1),
		Run:   noteUnpin,
	}
	moveCmd = &cobra.Command{
		Use:     "move id toSpace",
		Aliases: []string{"mv"},
		Short:   "Move note to other space",
		Args:    cobra.MinimumNArgs(2),
		Run:     noteMove,
	}
	spacesCmd = &cobra.Command{
		Use:     "spaces [id...]",
		Aliases: []string{"spc"},
		Short:   "Lists all or some spaces",
		Run:     noteSpaces,
	}
	importCmd = &cobra.Command{
		Use:     "import file [file...]",
		Aliases: []string{"imp"},
		Short:   "Import notes from JSON or YAML file",
		Args:    cobra.MinimumNArgs(1),
		Run:     noteImport,
	}
	exportCmd = &cobra.Command{
		Use:     "export",
		Aliases: []string{"exp"},
		Short:   "Export notes to JSON or YAML file",
		Run:     noteExport,
	}

	// Global arguments
	configPathArg  string
	storagePathArg string

	// Init argument
	forceArg bool

	// Add arguments
	fileArg   string
	pinnedArg bool

	// Remove arguments
	allInSpaceArg string
	noConfirmArg  bool
	permanentArg  bool

	// List arguments
	allArg        bool // also used in spaces
	sortByArg     string
	descendingArg bool
	limitArg      int
	offsetArg     int
	styleArg      string // also used in get
	colorArg      string // also used in get

	// Spaces arguments
	listArg bool

	// Import/Export arguments
	spacesArg     []string
	jsonArg       bool
	yamlArg       bool
	jsonIndentArg string
	jsonPrefixArg string
	yamlSpacesArg int
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
	globalFlags.StringVarP(&configPathArg, "config", "c", dfltConfig, "config file")
	globalFlags.StringVar(&storagePathArg, "db", dfltStore, "database store containing your notes")

	initFlags := initCmd.Flags()
	initFlags.BoolVar(&forceArg, "force", false, "determines if existing files will be overwritten")

	addFlags := addCmd.Flags()
	_ = addFlags.StringP("space", "s", DefaultSpace, "partitions the note into a space")
	addFlags.StringVarP(&fileArg, "file", "f", "", "the note is read from file")
	addFlags.BoolVarP(&pinnedArg, "pinned", "p", false, "pin your note to the top")

	removeFlags := removeCmd.Flags()
	removeFlags.StringVar(&allInSpaceArg, "all-in-space", "", "remove all notes in this space")
	removeFlags.BoolVar(&noConfirmArg, "no-confirm", false, "skip confirmation dialog")
	removeFlags.BoolVar(&permanentArg, "permanent", false, "note is completely removed from the db")

	sortKeys := getSortKeys()
	sortUsage := fmt.Sprintf("column to sort notes by (%v)", sortKeys)
	collectFlagSet := pflag.NewFlagSet("collect", pflag.ExitOnError)
	collectFlagSet.StringVarP(&sortByArg, "sort", "S", "id", sortUsage)
	collectFlagSet.IntVarP(&limitArg, "limit", "l", 0, "limit amount of notes listed, 0 means no limit")
	collectFlagSet.IntVarP(&offsetArg, "offset", "o", 0, "begin list notes at some offset (only if limit > 0)")
	collectFlagSet.BoolVarP(&descendingArg, "descending", "d", false, "descending order")

	printFlagSet := pflag.NewFlagSet("print", pflag.ExitOnError)
	printFlagSet.StringVar(&styleArg, "style", string(TitleStyle), "output style (plain, title)")
	printFlagSet.StringVar(&colorArg, "color", "auto", "color option (auto, no|never, yes|always)")

	listFlags := listCmd.Flags()
	listFlags.AddFlagSet(collectFlagSet)
	listFlags.AddFlagSet(printFlagSet)

	getFlags := getCmd.Flags()
	getFlags.AddFlagSet(printFlagSet)

	spacesFlags := spacesCmd.Flags()
	spacesFlags.BoolVarP(&allArg, "all", "a", false, "show hidden spaces")
	spacesFlags.BoolVarP(&listArg, "list", "l", false, "separate each space with a newline")
	spacesFlags.BoolVarP(&descendingArg, "descending", "d", false, "descending order")

	inoutFlagSet := pflag.NewFlagSet("inout", pflag.ExitOnError)
	inoutFlagSet.BoolVarP(&jsonArg, "json", "j", false, "JSON format")
	inoutFlagSet.BoolVarP(&yamlArg, "yaml", "y", false, "YAML format")

	importFlags := importCmd.Flags()
	importFlags.AddFlagSet(inoutFlagSet)
	importForceUsage := "all input files will use the format specified by either --json or --yaml"
	importFlags.BoolVarP(&listArg, "list", "l", false, "separate each id imported with a newline")
	importFlags.BoolVar(&forceArg, "force-format", false, importForceUsage)

	exportFlags := exportCmd.Flags()
	exportFlags.AddFlagSet(collectFlagSet)
	exportFlags.AddFlagSet(inoutFlagSet)
	exportFlags.StringSliceVarP(&spacesArg, "spaces", "s", []string{}, "limit export to notes from space(s)")
	exportFlags.BoolVar(&forceArg, "force", false, "determines if existing file will be overwritten")
	exportFlags.StringVarP(&jsonIndentArg, "indent", "i", "", "JSON indentation encoding option")
	exportFlags.StringVarP(&jsonPrefixArg, "prefix", "p", "", "JSON prefix encoding option")
	exportFlags.IntVarP(&yamlSpacesArg, "yaml-spaces", "P", 4, "YAML spaces encoding option")

	// These variables can exist in the config file or as environment variables as well
	viper.BindPFlag("db", globalFlags.Lookup("db"))
	viper.BindPFlag("add_space", addFlags.Lookup("space"))
	viper.BindPFlag("ls_spaces", listFlags.Lookup("spaces"))

	rootCmd.AddCommand(
		initCmd,
		addCmd, removeCmd,
		getCmd, listCmd, spacesCmd,
		editCmd, pinCmd, unpinCmd, moveCmd,
		importCmd, exportCmd,
	)
}

func quitError(loc string, err error) {
	stderr := fmt.Sprintf("error %v: %v", loc, err)
	fmt.Fprintln(os.Stderr, stderr)
	os.Exit(1)
}

func quit(msg string) {
	stderr := fmt.Sprintf("error: %v", msg)
	fmt.Fprintln(os.Stderr, stderr)
	os.Exit(1)
}
