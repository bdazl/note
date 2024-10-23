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
	"runtime"
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
	switch runtime.GOOS {
	case Linux:
		return defaultXdg(xdgConfigHome, defaultConfig)
	case Darwin:
		return darwinConfig()
	case Windows:
		return winLocalAppData()
	default:
		return "", fmt.Errorf("can't determine default configuration directory")
	}
}

func defaultDataDir() (string, error) {
	switch runtime.GOOS {
	case Linux:
		return defaultXdg(xdgDataHome, defaultData)
	case Darwin:
		return darwinCache()
	case Windows:
		return winAppData()
	default:
		return "", fmt.Errorf("can't determine default data directory")
	}
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

func darwinConfig() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir != "" {
		return configDir, nil
	}

	home, err := homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "Application Support"), nil
}

func darwinCache() (string, error) {
	// Under MacOS, we still follow freedesktop, but use more sensible default
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir != "" {
		return cacheDir, nil
	}

	home, err := homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "Caches"), nil
}

func winLocalAppData() (string, error) {
	localAppDir := os.Getenv("LOCALAPPDATA")
	if localAppDir == "" {
		userProfile, err := winUserProfile()
		if err != nil {
			return "", err
		}

		localAppDir = filepath.Join(userProfile, "AppData", "Local")
	}

	return localAppDir, nil
}

func winAppData() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		userProfile, err := winUserProfile()
		if err != nil {
			return "", err
		}

		appData = filepath.Join(userProfile, "AppData", "Roaming")
	}

	return appData, nil
}

func winUserProfile() (string, error) {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return "", fmt.Errorf("USERPROFILE not defined")
	}

	return userProfile, nil
}

func validFolder(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return stat.IsDir()
}
