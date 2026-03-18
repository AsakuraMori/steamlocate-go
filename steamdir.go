// Package steamlocate provides functionality for locating Steam installation directories
// and installed Steam applications.
//
// This is a Go port of the Rust steamlocate library.
//
// # Examples
//
// ## Locate the Steam installation and a specific game
//
//	steamDir, err := steamlocate.Locate()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println("Steam installation:", steamDir.Path())
//
//	const gmodAppID = 4000
//	gmod, library, err := steamDir.FindApp(gmodAppID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if gmod != nil {
//	    fmt.Printf("Found %s in %s\n", gmod.Name, library.Path())
//	}
//
// ## Get an overview of all libraries and apps on the system
//
//	libraries, err := steamDir.Libraries()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, library := range libraries {
//	    fmt.Println("Library:", library.Path())
//
//	    apps, err := library.Apps()
//	    if err != nil {
//	        continue
//	    }
//
//	    for _, app := range apps {
//	        fmt.Printf("  App %d - %s\n", app.AppID, app.Name)
//	    }
//	}
package steamlocate

import (
	"os"
	"path/filepath"
)

// SteamDir represents the Steam installation directory
type SteamDir struct {
	path string
}

// Locate attempts to find the Steam installation directory on the system
func Locate() (*SteamDir, error) {
	paths, err := locateSteamDir()
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, newLocateError("no steam installation found", nil)
	}

	return NewSteamDir(paths[0])
}

// LocateMultiple attempts to find all Steam installation directories on the system
// This is primarily useful on Linux where multiple installation methods may exist
func LocateMultiple() ([]*SteamDir, error) {
	paths, err := locateSteamDir()
	if err != nil {
		return nil, err
	}

	var steamDirs []*SteamDir
	for _, path := range paths {
		if sd, err := NewSteamDir(path); err == nil {
			steamDirs = append(steamDirs, sd)
		}
	}

	return steamDirs, nil
}

// NewSteamDir creates a SteamDir from a specific path
func NewSteamDir(path string) (*SteamDir, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, newMissingDirError()
		}
		return nil, newIOError(path, err)
	}

	if !info.IsDir() {
		return nil, newMissingDirError()
	}

	return &SteamDir{path: path}, nil
}

// Path returns the path to the Steam installation directory
func (sd *SteamDir) Path() string {
	return sd.path
}

// LibraryPaths returns the paths to all library folders
func (sd *SteamDir) LibraryPaths() ([]string, error) {
	libraryFoldersPath := filepath.Join(sd.path, "steamapps", "libraryfolders.vdf")
	return ParseLibraryFoldersFile(libraryFoldersPath)
}

// Libraries returns all Library instances for this Steam installation
func (sd *SteamDir) Libraries() ([]*Library, error) {
	paths, err := sd.LibraryPaths()
	if err != nil {
		// If libraryfolders.vdf doesn't exist or fails to parse,
		// try just the main steam directory
		if IsNotExist(err) {
			paths = []string{sd.path}
		} else {
			return nil, err
		}
	}

	// Always include the main steam directory
	found := false
	for _, p := range paths {
		if p == sd.path {
			found = true
			break
		}
	}
	if !found {
		paths = append([]string{sd.path}, paths...)
	}

	var libraries []*Library
	for _, path := range paths {
		if lib, err := NewLibraryFromDir(path); err == nil {
			libraries = append(libraries, lib)
		}
	}

	return libraries, nil
}

// FindApp searches through all libraries for a specific app by ID
// Returns (nil, nil, nil) if the app is not found
// Returns (app, library, nil) if found
// Returns (nil, nil, error) if an error occurred
func (sd *SteamDir) FindApp(appID uint32) (*App, *Library, error) {
	libraries, err := sd.Libraries()
	if err != nil {
		return nil, nil, err
	}

	for _, lib := range libraries {
		app, err := lib.GetApp(appID)
		if err != nil {
			continue
		}
		if app != nil {
			return app, lib, nil
		}
	}

	return nil, nil, nil
}

// CompatToolMapping returns the compatibility tool mapping (Proton configurations)
func (sd *SteamDir) CompatToolMapping() (CompatToolMapping, error) {
	return ParseCompatToolMapping(sd.path)
}

// Shortcuts returns all non-Steam game shortcuts
func (sd *SteamDir) Shortcuts() ([]*Shortcut, error) {
	return GetAllShortcuts(sd.path)
}

// SteamAppsPath returns the path to the steamapps directory
func (sd *SteamDir) SteamAppsPath() string {
	return filepath.Join(sd.path, "steamapps")
}

// UserDataPath returns the path to the userdata directory
func (sd *SteamDir) UserDataPath() string {
	return filepath.Join(sd.path, "userdata")
}

// ConfigPath returns the path to the config directory
func (sd *SteamDir) ConfigPath() string {
	return filepath.Join(sd.path, "config")
}
