package steamlocate

import (
	"os"
	"testing"
)

func TestUniverseFromInt(t *testing.T) {
	tests := []struct {
		input    int
		expected Universe
	}{
		{0, UniverseInvalid},
		{1, UniversePublic},
		{2, UniverseBeta},
		{3, UniverseInternal},
		{4, UniverseDev},
	}

	for _, tt := range tests {
		result := Universe(tt.input)
		if result != tt.expected {
			t.Errorf("Universe(%d) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestStateFlags(t *testing.T) {
	// Test StateFlags(0) returns Invalid
	sf := StateFlags(0)
	flags := sf.Flags()
	if len(flags) != 1 || flags[0] != StateFlagInvalid {
		t.Errorf("StateFlags(0).Flags() = %v, expected [StateFlagInvalid]", flags)
	}

	// Test individual flags
	sf = StateFlags(4) // 1 << 2 = FullyInstalled
	flags = sf.Flags()
	if len(flags) != 1 || flags[0] != StateFlagFullyInstalled {
		t.Errorf("StateFlags(4).Flags() = %v, expected [StateFlagFullyInstalled]", flags)
	}

	// Test multiple flags
	sf = StateFlags(6) // 1<<1 | 1<<2 = UpdateRequired | FullyInstalled
	flags = sf.Flags()
	if len(flags) != 2 {
		t.Errorf("StateFlags(6).Flags() length = %d, expected 2", len(flags))
	}
}

func TestParseTimestamp(t *testing.T) {
	// Test empty string
	if ts := parseTimestamp(""); ts != nil {
		t.Error("parseTimestamp(\"\") should return nil")
	}

	// Test "0"
	if ts := parseTimestamp("0"); ts != nil {
		t.Error("parseTimestamp(\"0\") should return nil")
	}

	// Test valid timestamp
	ts := parseTimestamp("1672176869")
	if ts == nil {
		t.Fatal("parseTimestamp(\"1672176869\") should return a time")
	}
	if ts.Unix() != 1672176869 {
		t.Errorf("Expected Unix timestamp 1672176869, got %d", ts.Unix())
	}
}

func TestParseAppManifest(t *testing.T) {
	manifest := `"AppState"
{
	"appid"		"4000"
	"installdir"		"GarrysMod"
	"name"		"Garry's Mod"
	"LastOwner"		"12345678901234567"
	"Universe"		"1"
	"StateFlags"		"6"
	"LastUpdated"		"1672176869"
	"buildid"		"8559806"
}`

	tmpDir := t.TempDir()
	manifestPath := tmpDir + "/appmanifest_4000.acf"

	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	app, err := ParseAppManifest(manifestPath)
	if err != nil {
		t.Fatalf("ParseAppManifest() error = %v", err)
	}

	if app.AppID != 4000 {
		t.Errorf("AppID = %d, expected 4000", app.AppID)
	}
	if app.InstallDir != "GarrysMod" {
		t.Errorf("InstallDir = %s, expected GarrysMod", app.InstallDir)
	}
	if app.Name != "Garry's Mod" {
		t.Errorf("Name = %s, expected Garry's Mod", app.Name)
	}
	if app.LastUser == nil || *app.LastUser != 12345678901234567 {
		t.Errorf("LastUser = %v, expected 12345678901234567", app.LastUser)
	}
	if app.Universe == nil || *app.Universe != UniversePublic {
		t.Errorf("Universe = %v, expected Public", app.Universe)
	}
	if app.StateFlags == nil || *app.StateFlags != 6 {
		t.Errorf("StateFlags = %v, expected 6", app.StateFlags)
	}
	if app.BuildID == nil || *app.BuildID != 8559806 {
		t.Errorf("BuildID = %v, expected 8559806", app.BuildID)
	}
}
