package steamlocate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLibraryFromDir(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	steamappsDir := filepath.Join(tmpDir, "steamapps")
	if err := os.MkdirAll(steamappsDir, 0755); err != nil {
		t.Fatalf("Failed to create steamapps dir: %v", err)
	}

	// Create some app manifest files
	manifests := []struct {
		name    string
		content string
	}{
		{
			name: "appmanifest_4000.acf",
			content: `"AppState"
{
	"appid"		"4000"
	"installdir"		"GarrysMod"
	"name"		"Garry's Mod"
}`,
		},
		{
			name: "appmanifest_230410.acf",
			content: `"AppState"
{
	"appid"		"230410"
	"installdir"		"Warframe"
	"name"		"Warframe"
}`,
		},
	}

	for _, m := range manifests {
		path := filepath.Join(steamappsDir, m.name)
		if err := os.WriteFile(path, []byte(m.content), 0644); err != nil {
			t.Fatalf("Failed to write manifest: %v", err)
		}
	}

	// Test library creation
	lib, err := NewLibraryFromDir(tmpDir)
	if err != nil {
		t.Fatalf("NewLibraryFromDir() error = %v", err)
	}

	if lib.Path() != tmpDir {
		t.Errorf("Path() = %s, expected %s", lib.Path(), tmpDir)
	}

	appIDs := lib.AppIDs()
	if len(appIDs) != 2 {
		t.Errorf("AppIDs() length = %d, expected 2", len(appIDs))
	}

	// Check that we can find the apps
	foundIDs := make(map[uint32]bool)
	for _, id := range appIDs {
		foundIDs[id] = true
	}

	if !foundIDs[4000] {
		t.Error("Expected to find app 4000")
	}
	if !foundIDs[230410] {
		t.Error("Expected to find app 230410")
	}
}

func TestLibraryGetApp(t *testing.T) {
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

	lib, err := NewLibraryFromDir(tmpDir)
	if err != nil {
		t.Fatalf("NewLibraryFromDir() error = %v", err)
	}

	// Test getting existing app
	app, err := lib.GetApp(4000)
	if err != nil {
		t.Fatalf("GetApp(4000) error = %v", err)
	}
	if app == nil {
		t.Fatal("GetApp(4000) returned nil")
	}
	if app.AppID != 4000 {
		t.Errorf("AppID = %d, expected 4000", app.AppID)
	}

	// Test getting non-existent app
	app, err = lib.GetApp(9999)
	if err != nil {
		t.Fatalf("GetApp(9999) error = %v", err)
	}
	if app != nil {
		t.Error("GetApp(9999) should return nil for non-existent app")
	}
}

func TestResolveAppDir(t *testing.T) {
	lib := &Library{path: "/home/user/Steam"}
	app := &App{InstallDir: "GarrysMod"}

	expected := filepath.Join("/home/user/Steam", "steamapps", "common", "GarrysMod")
	result := lib.ResolveAppDir(app)

	// Use filepath.Clean for cross-platform comparison
	if filepath.Clean(result) != filepath.Clean(expected) {
		t.Errorf("ResolveAppDir() = %s, expected %s", result, expected)
	}
}
