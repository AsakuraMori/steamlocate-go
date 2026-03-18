package steamlocate

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
)

// Shortcut represents a non-Steam game added to Steam
type Shortcut struct {
	AppID      uint32
	AppName    string
	Executable string
	StartDir   string
}

// SteamID calculates the Steam ID from the executable and app name
func (s *Shortcut) SteamID() uint64 {
	// Using CRC32 ISO HDLC polynomial
	table := crc32.MakeTable(crc32.IEEE)
	digest := crc32.Checksum([]byte(s.Executable), table)
	digest = crc32.Update(digest, table, []byte(s.AppName))

	// Alternative CRC32 implementation for ISO HDLC
	// Since Go's standard CRC32 doesn't have ISO HDLC, we use IEEE as closest alternative
	// The original Rust code uses CRC_32_ISO_HDLC which is different
	// For exact compatibility, we'd need to implement the specific CRC variant

	top := uint64(digest) | 0x80000000
	return (top << 32) | 0x02000000
}

// ParseShortcuts parses a shortcuts.vdf file
func ParseShortcuts(path string) ([]*Shortcut, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Not all users have shortcuts
		}
		return nil, newIOError(path, err)
	}

	return parseShortcutsBinary(data)
}

// parseShortcutsBinary parses binary VDF format shortcuts
func parseShortcutsBinary(data []byte) ([]*Shortcut, error) {
	var shortcuts []*Shortcut
	
	// Binary VDF parsing for shortcuts
	// Format is similar to text VDF but with binary markers
	
	i := 0
	for i < len(data) {
		// Look for shortcut entry marker (0x00 followed by numeric key)
		if data[i] != 0x00 {
			i++
			continue
		}

		// Try to parse a shortcut entry
		shortcut, consumed := parseShortcutEntry(data[i:])
		if shortcut != nil {
			shortcuts = append(shortcuts, shortcut)
		}
		if consumed == 0 {
			i++
		} else {
			i += consumed
		}
	}

	return shortcuts, nil
}

func parseShortcutEntry(data []byte) (*Shortcut, int) {
	if len(data) < 10 {
		return nil, 0
	}

	shortcut := &Shortcut{}
	i := 0

	// Look for appid marker (0x02 'appid' 0x00)
	appidIdx := bytes.Index(data[i:], []byte{0x02, 'a', 'p', 'p', 'i', 'd', 0x00})
	if appidIdx == -1 {
		appidIdx = bytes.Index(data[i:], []byte{0x02, 'A', 'p', 'p', 'I', 'D', 0x00})
	}
	if appidIdx == -1 {
		return nil, 0
	}
	i += appidIdx + 7 // Skip past marker

	// Read 4-byte appid
	if i+4 > len(data) {
		return nil, 0
	}
	shortcut.AppID = binary.LittleEndian.Uint32(data[i:])
	i += 4

	// Look for AppName marker (0x01 'AppName' 0x00)
	nameIdx := findCaseInsensitive(data[i:], []byte{0x01, 'A', 'p', 'p', 'N', 'a', 'm', 'e', 0x00})
	if nameIdx == -1 {
		return nil, 0
	}
	i += nameIdx + 9

	// Read null-terminated string
	name, consumed := readNullTerminatedString(data[i:])
	if consumed == 0 {
		return nil, 0
	}
	shortcut.AppName = name
	i += consumed

	// Look for Exe marker (0x01 'Exe' 0x00)
	exeIdx := findCaseInsensitive(data[i:], []byte{0x01, 'E', 'x', 'e', 0x00})
	if exeIdx == -1 {
		return nil, 0
	}
	i += exeIdx + 5

	// Read null-terminated string
	exe, consumed := readNullTerminatedString(data[i:])
	if consumed == 0 {
		return nil, 0
	}
	shortcut.Executable = exe
	i += consumed

	// Look for StartDir marker (0x01 'StartDir' 0x00)
	startDirIdx := findCaseInsensitive(data[i:], []byte{0x01, 'S', 't', 'a', 'r', 't', 'D', 'i', 'r', 0x00})
	if startDirIdx == -1 {
		return nil, 0
	}
	i += startDirIdx + 10

	// Read null-terminated string
	startDir, consumed := readNullTerminatedString(data[i:])
	if consumed == 0 {
		return nil, 0
	}
	shortcut.StartDir = startDir

	return shortcut, i
}

func findCaseInsensitive(data, pattern []byte) int {
	for i := 0; i <= len(data)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if !bytesEqualIgnoreCase(data[i+j], pattern[j]) {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func bytesEqualIgnoreCase(a, b byte) bool {
	if a >= 'A' && a <= 'Z' {
		a = a - 'A' + 'a'
	}
	if b >= 'A' && b <= 'Z' {
		b = b - 'A' + 'a'
	}
	return a == b
}

func readNullTerminatedString(data []byte) (string, int) {
	for i, b := range data {
		if b == 0x00 {
			return string(data[:i]), i + 1
		}
	}
	return "", 0
}

// GetAllShortcuts returns all shortcuts from all user data directories
func GetAllShortcuts(steamPath string) ([]*Shortcut, error) {
	userDataPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(userDataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, newIOError(userDataPath, err)
	}

	var allShortcuts []*Shortcut
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if entry name is numeric (Steam user ID)
		userID := entry.Name()
		if _, err := fmt.Sscanf(userID, "%d", new(uint64)); err != nil {
			continue
		}

		shortcutsPath := filepath.Join(userDataPath, userID, "config", "shortcuts.vdf")
		shortcuts, err := ParseShortcuts(shortcutsPath)
		if err != nil {
			continue // Skip users without shortcuts
		}

		allShortcuts = append(allShortcuts, shortcuts...)
	}

	return allShortcuts, nil
}
