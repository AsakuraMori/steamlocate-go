package steamlocate

import (
	"strings"
	"testing"
)

func TestParseVDF(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple key-value",
			input: `"Root"
{
	"key1" "value1"
	"key2" "value2"
}`,
			wantErr: false,
		},
		{
			name: "nested structure",
			input: `"Root"
{
	"child"
	{
		"grandchild" "value"
	}
}`,
			wantErr: false,
		},
		{
			name: "libraryfolders example",
			input: `"libraryfolders"
{
	"0"
	{
		"path"		"/home/user/.local/share/Steam"
		"label"		""
		"contentid"		"123456789"
		"totalsize"		"0"
		"update_clean_bytes_tally"		"0"
		"time_last_update_corruption"		"0"
		"apps"
		{
			"4000"		"123456789"
		}
	}
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseVDF(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVDF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if node == nil {
				t.Error("ParseVDF() returned nil node")
			}
		})
	}
}

func TestParseLibraryFolders(t *testing.T) {
	input := `"libraryfolders"
{
	"0"
	{
		"path"		"/home/user/.local/share/Steam"
		"apps"
		{
			"4000"		"123456789"
		}
	}
	"1"
	{
		"path"		"/mnt/games/Steam"
		"apps"
		{
			"230410"		"987654321"
		}
	}
}`

	paths, err := ParseLibraryFolders(input)
	if err != nil {
		t.Fatalf("ParseLibraryFolders() error = %v", err)
	}

	if len(paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(paths))
	}

	expected := []string{
		"/home/user/.local/share/Steam",
		"/mnt/games/Steam",
	}

	for i, exp := range expected {
		if i >= len(paths) {
			break
		}
		if !strings.Contains(paths[i], exp) && !strings.Contains(exp, paths[i]) {
			t.Errorf("Path %d: expected %s, got %s", i, exp, paths[i])
		}
	}
}

func TestVDFNodeGetters(t *testing.T) {
	input := `"Root"
{
	"string" "hello"
	"number" "123"
	"nested"
	{
		"value" "world"
	}
}`

	root, err := ParseVDF(input)
	if err != nil {
		t.Fatalf("ParseVDF() error = %v", err)
	}

	rootNode := root.Get("Root")
	if rootNode == nil {
		t.Fatal("Root node is nil")
	}

	if str := rootNode.Get("string"); str == nil || str.GetString() != "hello" {
		t.Errorf("Expected 'hello', got %v", str)
	}

	if num := rootNode.Get("number"); num == nil || num.GetInt() != 123 {
		t.Errorf("Expected 123, got %v", num)
	}

	if nested := rootNode.Get("nested"); nested == nil {
		t.Error("Nested node is nil")
	} else {
		if val := nested.Get("value"); val == nil || val.GetString() != "world" {
			t.Errorf("Expected 'world', got %v", val)
		}
	}
}
