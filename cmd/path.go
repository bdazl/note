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
	"os/user"
	"path/filepath"
)

const (
	note = "note"

	defaultConfigName  = "note.yaml"
	defaultStorageName = "note.db"

	defaultConfig = ".config/note"
	defaultData   = ".local/share/note"

	xdgConfigHome = "XDG_CONFIG_HOME"
	xdgDataHome   = "XDG_DATA_HOME"
)

// The user home directory
// If this function does not return an error, the path is a valid folder
// $HOME variable takes precedence.
func homeDir() (string, error) {
	const (
		preError    = "error determining home directory"
		noValidHome = "neither $HOME nor current user home dir is a valid directory"
	)
	envHome := os.Getenv("HOME")

	// If $HOME is defined and valid, that takes precedence
	if validFolder(envHome) {
		return envHome, nil
	}

	// If $HOME is not defined, look at user settings
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("%v: %w", preError, err)
	}

	if !validFolder(user.HomeDir) {
		return "", fmt.Errorf("%v: %v", preError, noValidHome)
	}
	return user.HomeDir, nil
}

func defaultConfigPath() (string, error) {
	cfgDir, err := defaultConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cfgDir, note, "note.yaml"), nil
}

func defaultStoragePath() (string, error) {
	dataDir, err := defaultDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, note, defaultStorageName), nil
}

// The default directory of configuration files
func defaultConfigDir() (string, error) {
	return defaultXdg(xdgConfigHome, defaultConfig)
}

func defaultDataDir() (string, error) {
	return defaultXdg(xdgDataHome, defaultData)
}

func defaultXdg(env, rel string) (string, error) {
	envVal := os.Getenv(env)

	// $XDG_* takes precedence
	if validFolder(envVal) {
		return envVal, nil
	}

	home, err := homeDir()
	if err != nil {
		return "", err
	}

	// If we got home, we know it's a valid directory
	// The underlying directory must not necessarily exist yet
	return filepath.Join(home, rel), nil
}

func validFolder(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return stat.IsDir()
}
