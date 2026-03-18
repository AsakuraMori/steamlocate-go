package steamlocate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSteamDir(t *testing.T) {
	// Test with valid directory
	tmpDir := t.TempDir()
	sd, err := NewSteamDir(tmpDir)
	if err != nil {
		t.Fatalf("NewSteamDir() error = %v", err)
	}
	if sd.Path() != tmpDir {
		t.Errorf("Path() = %s, expected %s", sd.Path(), tmpDir)
	}

	// Test with non-existent directory
	_, err = NewSteamDir("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("NewSteamDir() should return error for non-existent directory")
	}

	// Test with file instead of directory
	tmpFile := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_, err = NewSteamDir(tmpFile)
	if err == nil {
		t.Error("NewSteamDir() should return error for file (not directory)")
	}
}

func TestSteamDirLibraryPaths(t *testing.T) {
	// Create a mock Steam directory structure
	tmpDir := t.TempDir()
	steamappsDir := filepath.Join(tmpDir, "steamapps")
	if err := os.MkdirAll(steamappsDir, 0755); err != nil {
		t.Fatalf("Failed to create steamapps dir: %v", err)
	}

	libraryfolders := `"libraryfolders"
{
	"0"
	{
		"path"		"` + filepath.ToSlash(tmpDir) + `"
		"apps"
		{
		}
	}
	"1"
	{
		"path"		"/mnt/games/SteamLibrary"
		"apps"
		{
		}
	}
}`

	libraryfoldersPath := filepath.Join(steamappsDir, "libraryfolders.vdf")
	if err := os.WriteFile(libraryfoldersPath, []byte(libraryfolders), 0644); err != nil {
		t.Fatalf("Failed to write libraryfolders.vdf: %v", err)
	}

	sd := &SteamDir{path: tmpDir}
	paths, err := sd.LibraryPaths()
	if err != nil {
		t.Fatalf("LibraryPaths() error = %v", err)
	}

	if len(paths) != 2 {
		t.Errorf("LibraryPaths() returned %d paths, expected 2", len(paths))
	}
}

func TestSteamDirFindApp(t *testing.T) {
	// Create a mock Steam directory with one app
	tmpDir := t.TempDir()
	steamappsDir := filepath.Join(tmpDir, "steamapps")
	if err := os.MkdirAll(steamappsDir, 0755); err != nil {
		t.Fatalf("Failed to create steamapps dir: %v", err)
	}

	manifest := `"AppState"
{
	"appid"		"4000"
	"installdir"		"GarrysMod"
	"name"		"Garry's Mod"
}`

	manifestPath := filepath.Join(steamappsDir, "appmanifest_4000.acf")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Create libraryfolders.vdf
	libraryfolders := `"libraryfolders"
{
	"0"
	{
		"path"		"` + filepath.ToSlash(tmpDir) + `"
		"apps"
		{
			"4000"		"123456789"
		}
	}
}`

	libraryfoldersPath := filepath.Join(steamappsDir, "libraryfolders.vdf")
	if err := os.WriteFile(libraryfoldersPath, []byte(libraryfolders), 0644); err != nil {
		t.Fatalf("Failed to write libraryfolders.vdf: %v", err)
	}

	sd := &SteamDir{path: tmpDir}

	// Test finding existing app
	app, library, err := sd.FindApp(4000)
	if err != nil {
		t.Fatalf("FindApp(4000) error = %v", err)
	}
	if app == nil {
		t.Fatal("FindApp(4000) returned nil app")
	}
	if library == nil {
		t.Fatal("FindApp(4000) returned nil library")
	}
	if app.AppID != 4000 {
		t.Errorf("AppID = %d, expected 4000", app.AppID)
	}

	// Test finding non-existent app
	app, library, err = sd.FindApp(9999)
	if err != nil {
		t.Fatalf("FindApp(9999) error = %v", err)
	}
	if app != nil || library != nil {
		t.Error("FindApp(9999) should return nil for non-existent app")
	}
}
