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

	"github.com/spf13/viper"
)

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
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

	viper.SetDefault("db", dfltStore)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		cmdPath := currentCmd.CommandPath()
		// Only print warning message when we the user is not explicitly trying to create this file
		if cmdPath != "note init" {
			fmt.Fprintln(os.Stderr, "Could not read config file, consider running: note init")
		}
	}
}

func dbFilename() string {
	return viper.GetString("db")
}
