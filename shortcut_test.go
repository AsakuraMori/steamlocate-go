package steamlocate

import (
	"testing"
)

func TestShortcutSteamID(t *testing.T) {
	// The Steam ID calculation depends on CRC32
	// This is a basic test to ensure it doesn't panic
	shortcut := &Shortcut{
		AppID:      2786274309,
		AppName:    "Anki",
		Executable: "\"anki\"",
		StartDir:   "\"./\"",
	}

	steamID := shortcut.SteamID()
	if steamID == 0 {
		t.Error("SteamID() returned 0")
	}

	// The Steam ID should have the high bit set and the magic number
	if steamID&0x02000000 == 0 {
		t.Error("SteamID() missing magic number")
	}
}

func TestFindCaseInsensitive(t *testing.T) {
	tests := []struct {
		data     []byte
		pattern  []byte
		expected int
	}{
		{
			data:     []byte{0x01, 'A', 'p', 'p', 'N', 'a', 'm', 'e', 0x00},
			pattern:  []byte{0x01, 'a', 'p', 'p', 'n', 'a', 'm', 'e', 0x00},
			expected: 0,
		},
		{
			data:     []byte{0x00, 0x01, 'E', 'x', 'e', 0x00},
			pattern:  []byte{0x01, 'e', 'x', 'e', 0x00},
			expected: 1,
		},
		{
			data:     []byte{0x00, 0x00, 0x00},
			pattern:  []byte{0x01, 't', 'e', 's', 't', 0x00},
			expected: -1,
		},
	}

	for _, tt := range tests {
		result := findCaseInsensitive(tt.data, tt.pattern)
		if result != tt.expected {
			t.Errorf("findCaseInsensitive(%v, %v) = %d, expected %d", tt.data, tt.pattern, result, tt.expected)
		}
	}
}

func TestReadNullTerminatedString(t *testing.T) {
	tests := []struct {
		data     []byte
		expected string
		consumed int
	}{
		{
			data:     []byte{'h', 'e', 'l', 'l', 'o', 0x00, 'w', 'o', 'r', 'l', 'd'},
			expected: "hello",
			consumed: 6,
		},
		{
			data:     []byte{0x00},
			expected: "",
			consumed: 1,
		},
		{
			data:     []byte{'t', 'e', 's', 't'},
			expected: "",
			consumed: 0,
		},
	}

	for _, tt := range tests {
		str, consumed := readNullTerminatedString(tt.data)
		if str != tt.expected || consumed != tt.consumed {
			t.Errorf("readNullTerminatedString(%v) = (%s, %d), expected (%s, %d)",
				tt.data, str, consumed, tt.expected, tt.consumed)
		}
	}
}
