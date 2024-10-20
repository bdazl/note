# note

No fuzz terminal application for quick and simple note taking. 

## Installation

This is a [Go](https://go.dev) application, and currently there are no pre-packaged binaries.
To install this program you need the following pre-requisites:
1. [Go](https://go.dev/doc/install)
2. [go-sqlite3](https://github.com/mattn/go-sqlite3?tab=readme-ov-file#installation)

After these are installed, the installation of this program is done by running
```bash
go install github.com/bdazl/note@latest
```

## Usage

`note` uses [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper)
to handle command line arguments and configuration. This means that any sub-command that you use
can be invoked with the `-h` flag, to get information about usage and parameter to that command.

```
Available Commands:
  add         Add new note
  completion  Generate the autocompletion script for the specified shell
  edit        Edit content of note
  export      Export notes to JSON or YAML file
  get         Get specific note(s)
  help        Help about any command
  id          Lists all or some IDs
  import      Import notes from JSON or YAML file
  init        Initialize note configuration and database
  list        Lists notes from one or more spaces
  move        Move note to another space
  pin         Pin note(s) to top
  remove      Remove note(s) with id(s)
  space       Lists all or some spaces
  table       Lists available notes in a table format
  unpin       Unpin note(s) from top
```

### Init note

The first time you use note you need to initialize a configuration and a database. This is done by
calling the `init` sub-command:
```bash
note init
```

This will store a configuration and a database file. The locations of these files can
be modified by specifying `-c` and `--db` respectively. These arguments are global to all commands
and, if specified, will override the config initialized above. The configuration can be a
[wide array of formats](https://github.com/spf13/viper?tab=readme-ov-file#what-is-viper).

By default the [Freedesktop XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/)
is used. Mainly `note` looks at the `$XDG_DATA_HOME` and `$XDG_CONFIG_HOME` to determine default directories.
If these variables are not defined, default values are set according to specification.

### Add note

There are three main methods of adding your note. The quickest way is to simply invoke note like this:
```bash
note add This is My First Note
note add "This is Another Note"
```

`note` will respond with an id that you use if you want to modify or access it later:
```bash
Created note: 1
```

You can write a note with your default `$EDITOR`:
```bash
note add
```

If no such editor exist, a default has been chosen for you. It is possible to define the editor in the
configuration file.

A note can also be created by specifying a file that you want note to read:
```bash
note add -f some.txt
```
`note` can also read from standard input (`-`):
```bash
echo note made by other program | note add -f -
```

### Spaces

All notes belong to a space, which can be any string you define:
```bash
note add -s MySpace This is a Song
```

By default, your notes will fall into the `main` space. You can list all spaces occupied with one or
more notes, by running:
```bash
note space [id...]
```

If you specify one (or more) id:s in the above command, only spaces occupied by the notes you specify will
be shown. You can also use the alias `note spc [ids...]` which is equivalent to the previous statement.

To list IDs occupied by a space, you can use the following command, and similarly to the `space` command
above, if you specify one or more positional arguments - only ID's in those spaces will be shown
```bash
note id [space...]
```

### Content of notes

To get an overview of your notes, it's often useful to get a table. This can be done simply with:
```bash
note
```

If you want to have control over how many preview words will be shown, use the specific sub-command:
```bash
note table --preview 5 [space...]
```

To retrieve the content of specific note(s), invoke `get`:
```bash
note get id [id...]
```

To list the full content of your notes, invoke the `list` (or `ls`) command, you can specify any number
of space(s) as a filter. There are also a number of sorting, limiting and styling options to this command:
```bash
note ls [space...]
```

### Edit notes

Edit note in `$EDITOR`:
```bash
note edit id
```

Pin/unpin note(s) to the top:
```bash
note pin id [id...]
note unpin id [id...]
```

Move note to space:
```bash
note move space id [id...]
```

## Import/export

Notes can be exported and imported to [JSON](https://en.wikipedia.org/wiki/JSON) or
[YAML](https://en.wikipedia.org/wiki/YAML). By default the export is printed to standard output:
```bash
note export [file]
```

## Configuration

At the moment the configuration file is not that useful, this is in `TODO`-stage at the moment. Use with
care.

## License

MIT License
