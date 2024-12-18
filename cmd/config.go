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
	"runtime"

	"github.com/bdazl/note/db"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const (
	Linux   = "linux"
	Darwin  = "darwin"
	Windows = "windows"

	ViperDb     = "db"
	ViperEditor = "editor"
	ViperSpace  = "space"
	ViperStyle  = "style"
	ViperColor  = "color"

	DefaultSpace = "main"
)

func initConfig() {
	if configPathArg != "" {
		viper.SetConfigFile(configPathArg)
	} else {
		cfgDir, err := defaultConfigDir()
		if err != nil {
			quitError("init config dir", err)
		}

		viper.AddConfigPath(cfgDir)
		viper.SetConfigName(note)
		viper.SetConfigType("yaml")
	}

	dfltStore, err := defaultStoragePath()
	if err != nil {
		quitError("init default store", err)
	}

	viper.SetDefault(ViperDb, dfltStore)
	viper.SetDefault(ViperEditor, defaultEditor())
	viper.SetDefault(ViperSpace, DefaultSpace)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		cmdPath := currentCmd.CommandPath()

		// Only print warning message when the config should be in place
		if cmdPathWantsConfig(cmdPath) {
			fmt.Fprintln(os.Stderr, "WARNING: Could not read config file, consider running: note init")
		}
	}
}

func cmdPathWantsConfig(cmdPath string) bool {
	switch cmdPath {
	case "note init", "note version":
		return false
	default:
		return true
	}
}

func defaultEditor() string {
	switch runtime.GOOS {
	case Linux:
		return "nano" // it pains me to put this here...
	case Darwin:
		return "open -a TextEdit"
	case Windows:
		return "notepad"
	default:
		return "vim"
	}
}

func dbFilename() string {
	return viper.GetString(ViperDb)
}

func dbOpen() *db.DB {
	d, err := db.Open(dbFilename())
	if err != nil {
		quitError("db open", err)
	}
	return d
}

func styleColorOpts() (Style, bool, error) {
	var (
		style   Style
		doColor bool
	)

	styleStr := viper.GetString(ViperStyle)
	stylized := Style(styleStr)
	switch stylized {
	case MinimalStyle, LightStyle, FullStyle:
		style = stylized
	default:
		return "", false, fmt.Errorf("unrecognized style")
	}

	colorStr := viper.GetString(ViperColor)
	switch colorStr {
	case "auto":
		doColor = !color.NoColor
	case "yes", "always":
		doColor = true
	case "no", "never":
		doColor = false
	}

	return style, doColor, nil
}
