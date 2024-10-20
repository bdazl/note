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
		Run:   noteTable,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			currentCmd = cmd
			initConfig()
		},
		Long: `A command line note taking app, to store your short form notes.

The program is designed to be quickly used to jot down text into notes, similar
to a bulletin board. Hence, notes does not have names, but rather they are
identified by an ID. Some organization is often desired, however, and so a note
is also placed in a so-called space. A note can only be in one space at a time,
but it's possible to move them them.

To start using note, run the 'note init' command, to initialize a configuration
and database file.

Notes can be pinned, which means that they will always be placed at the top -
or bottom - depending on the sort order. Pin or unpin with 'note pin' and
'note unpin' respectively.

Running note without any arguments is the same as the sub command 'note table',
listing all your notes in a table format. Some sub commands has short form
aliases. Like 'note ls', which is short for 'note list', or 'note rm' - short for
'note remove'. For information about specific sub commands, use the '--help' or
'-h' option. For example: 'note export -h'.`,
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize note configuration and database",
		Run:   noteInit,
		Long: `Use this to create an initial configuration and database file.

The configuration file holds default values, to make some operations less tedious.
All values in the configuration file can be overridden by command line arguments.

The database is a sqlite3 file with a very simple schema. This is used to store
your notes and is used in every operation. The location of this file can at any
time be altered in your configuration file. If you want to re-create the database,
you must first remove (or change location) of the file, and then run 'note init'
again.`,
	}
	addCmd = &cobra.Command{
		Use:   "add [note...]",
		Short: "Add new note",
		Run:   noteAdd,
		Long: `The add sub-command has three distinct ways to create new notes.

The simplest form of this command is to simply run 'note add', without arguments.
This will open up the editor of your choice; as specified in either the
configuration file, by command line argument, or by the environemnt variable
$EDITOR. Empty notes will be rejected.

The second way to add new notes is to word them out as argument(s), like so:
'note add These Arguments Becomes Content'. The arguments will be joined with
spaces, to resemble how they were input: "These Arguments Becomes Content".

The final way to create notes is by specifying a file, that will be read by note.
This file can be a text file or the special character '-', indicating standard
input.`,
	}
	removeCmd = &cobra.Command{
		Use:     "remove id [id...]",
		Aliases: []string{"rm", "del"},
		Short:   "Remove note(s) with id(s)",
		Run:     noteRemove,
		Long: `Remove one or many notes, by their respective ID's.

Removal is by default an operation that moves the notes to the space '.trash'.
To remove notes permanently you need to specify the '--permanent' flag. It is
possible to remove all notes in a space, by specifying the '--all-in-space'
argument, followed by the space you want to empty.`,
	}
	getCmd = &cobra.Command{
		Use:   "get id [id...]",
		Short: "Get specific note(s)",
		Args:  cobra.MinimumNArgs(1),
		Run:   noteGet,
		Long: `Print the contents of one or more note ID's.

The order of the notes will be the same as the input order.
For style and coloring options, see 'note list -h'.`,
	}
	listCmd = &cobra.Command{
		Use:     "list [space...]",
		Aliases: []string{"ls"},
		Short:   "Lists notes from one or more spaces",
		Run:     noteList,
		Long: `All the notes content in the input spaces will be output.

If no spaces are given, all your notes will be printed.

Sort options:
The --sort, or -S option determines the main sort column. The pinned
notes will always be the first when the order is ascending and last
when the order is descending. Reverse the print order by using the -d or
--descending argument. Available sort columns are:
* id (default)
* created (time)
* updated (last updated time)
* space
* content

The limit and offset options can be used to limit the amount of printed
notes and the offset determines the starting note to print. This can be
used to paginate the output.

Style options:
There are at present two styles: raw, light and full.
The raw style is meant to be showing only the most essential output.
When used with the 'note list' sub command, it will only print the content
of your notes. The light option will show you some context and the full
option will show everything.

Color options: auto, no or never, yes or always.
auto means that note will default to colors, if stdout is not connected
with a pipe or similar. To force color use the always, or equivalently
the yes option. The option never, or equivalently no, means never show
text in color.`,
	}
	tableCmd = &cobra.Command{
		Use:     "table [space...]",
		Aliases: []string{"tbl"},
		Short:   "Lists available notes in a table format",
		Run:     noteTable,
		Long: `Print a table of notes with their properties

If no spaces are input, notes from all spaces will be included.

The --preview, or -p option is an integera count; used to determine how
many preview words will be shown of the content in the notes.
If 0 is chosen, preview is disabled. If the note is in binary format
a word is defined as 5 characters.

Sort options can be found by running: 'note list -h'`,
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
		Use:     "move space id [id...]",
		Aliases: []string{"mv"},
		Short:   "Move note to another space",
		Args:    cobra.MinimumNArgs(2),
		Run:     noteMove,
	}
	idCmd = &cobra.Command{
		Use:     "id [space...]",
		Aliases: []string{"ids"},
		Short:   "Lists all or some IDs",
		Run:     noteId,
		Long: `Print the available IDs of notes.

If no spaces are given, all ID's will be listed.
By specifying one or more spaces, only the ID's occupied by those notes
will be shown.`,
	}
	spaceCmd = &cobra.Command{
		Use:     "space [id...]",
		Aliases: []string{"spaces", "spc"},
		Short:   "Lists all or some spaces",
		Run:     noteSpace,
		Long: `Print available spaces occupied by notes.

If no ID's are given, all spaces will be printed.
By specifying ID's of notes, only the spaces occupied by those notes
will be shown.`,
	}
	importCmd = &cobra.Command{
		Use:     "import file [file...]",
		Aliases: []string{"imp"},
		Short:   "Import notes from JSON or YAML file",
		Args:    cobra.MinimumNArgs(1),
		Run:     noteImport,
		Long: `Import many notes from a JSON or YAML file.

The top level is a list and each item is an object containg the following fields:
* content - string (required)
* space - string (optional; if not specified)
* created - date string (optional; if not specified, current time is chosen)
* last_updated - date string (optional; if not specified, current time is chosen)
* pinned - bool (optional; default: false)

Files will only be imported once (per run), no checks for duplicate notes are made.`,
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

	// Table
	previewArg uint

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

	tableFlags := tableCmd.Flags()
	tableFlags.AddFlagSet(collectFlagSet)
	tableFlags.UintVarP(&previewArg, "preview", "p", 5, "preview word count to display in table")

	getFlags := getCmd.Flags()
	getFlags.AddFlagSet(printFlagSet)

	idFlags := idCmd.Flags()
	idFlags.BoolVarP(&listArg, "list", "l", false, "separate each ID with a newline")
	idFlags.BoolVarP(&descendingArg, "descending", "d", false, "descending order")

	spaceFlags := spaceCmd.Flags()
	spaceFlags.BoolVarP(&listArg, "list", "l", false, "separate each space with a newline")
	spaceFlags.BoolVarP(&descendingArg, "descending", "d", false, "descending order")

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
	viper.BindPFlag(ViperDb, globalFlags.Lookup("db"))
	viper.BindPFlag(ViperAddSpace, addFlags.Lookup("space"))

	rootCmd.AddCommand(
		initCmd,
		addCmd, removeCmd,
		getCmd, listCmd, tableCmd, idCmd, spaceCmd,
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
