package steamlocate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Library represents a Steam library folder containing installed apps
type Library struct {
	path   string
	appIDs []uint32
}

// NewLibraryFromDir creates a Library from a directory path
func NewLibraryFromDir(path string) (*Library, error) {
	// Read the steamapps directory to get installed apps
	steamappsPath := filepath.Join(path, "steamapps")
	entries, err := os.ReadDir(steamappsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Try to create the library anyway, it might just be empty
			return &Library{
				path:   path,
				appIDs: []uint32{},
			}, nil
		}
		return nil, newIOError(steamappsPath, err)
	}

	var appIDs []uint32
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, "appmanifest_") && strings.HasSuffix(name, ".acf") {
			appIDStr := strings.TrimPrefix(name, "appmanifest_")
			appIDStr = strings.TrimSuffix(appIDStr, ".acf")
			if appID, err := strconv.ParseUint(appIDStr, 10, 32); err == nil {
				appIDs = append(appIDs, uint32(appID))
			}
		}
	}

	return &Library{
		path:   path,
		appIDs: appIDs,
	}, nil
}

// Path returns the library's installation directory
func (l *Library) Path() string {
	return l.path
}

// AppIDs returns the list of Application IDs in this library
func (l *Library) AppIDs() []uint32 {
	return l.appIDs
}

// GetApp attempts to get an App by its ID from this library
// Returns (nil, nil) if the app is not in this library
// Returns (nil, error) if there was an error reading the manifest
func (l *Library) GetApp(appID uint32) (*App, error) {
	found := false
	for _, id := range l.appIDs {
		if id == appID {
			found = true
			break
		}
	}
	if !found {
		return nil, nil
	}

	manifestPath := filepath.Join(l.path, "steamapps", fmt.Sprintf("appmanifest_%d.acf", appID))
	return ParseAppManifest(manifestPath)
}

// Apps returns all apps in this library
func (l *Library) Apps() ([]*App, error) {
	var apps []*App
	for _, appID := range l.appIDs {
		app, err := l.GetApp(appID)
		if err != nil {
			continue // Skip apps that fail to parse
		}
		if app != nil {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

// ResolveAppDir returns the installation directory for an app in this library
func (l *Library) ResolveAppDir(app *App) string {
	return ResolveAppDir(l.path, app)
}

// ParseLibraryFoldersFile parses the libraryfolders.vdf file at the given path
func ParseLibraryFoldersFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, newMissingFileError(ParseErrorKindLibraryFolders, path)
		}
		return nil, newIOError(path, err)
	}

	paths, err := ParseLibraryFolders(string(data))
	if err != nil {
		return nil, newParseError(ParseErrorKindLibraryFolders, path, err)
	}

	return paths, nil
}
