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

func noteInit(cmd *cobra.Command, args []string) {
	forceInform := false
	dbF, err := filepath.Abs(storagePathArg) // When doing init we explicitly want the command line option
	if err != nil {
		quitError("db path", err)
	}

	// We force viper to set value to the (abs path of the) command line option here
	// This is because init may be re-ran after we have a valid config setup, and we don't want to then source
	// the option from the config file.
	viper.Set("db", dbF)

	mkdir(filepath.Dir(configPathArg))
	mkdir(filepath.Dir(dbF))

	if !forceArg && exists(configPathArg) {
		fmt.Fprintln(os.Stderr, "Config file already exists")
		forceInform = true
	} else {
		fmt.Printf("Writing config file: %v\n", configPathArg)
		err := viper.WriteConfig()
		if err != nil {
			quitError("writing config", err)
		}
	}

	if !forceArg && exists(dbF) {
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
