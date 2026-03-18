//go:build !windows && !darwin

package steamlocate

import (
	"os"
	"path/filepath"
)

// locateSteamDir returns possible Steam installation directories on Linux
func locateSteamDir() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, newNoHomeError()
	}

	snapDir := os.Getenv("SNAP_USER_DATA")
	if snapDir == "" {
		snapDir = filepath.Join(homeDir, "snap")
	}

	paths := []string{
		// Flatpak steam install directories
		filepath.Join(homeDir, ".var/app/com.valvesoftware.Steam/.local/share/Steam"),
		filepath.Join(homeDir, ".var/app/com.valvesoftware.Steam/.steam/steam"),
		filepath.Join(homeDir, ".var/app/com.valvesoftware.Steam/.steam/root"),
		// Standard install directories
		filepath.Join(homeDir, ".local/share/Steam"),
		filepath.Join(homeDir, ".steam/steam"),
		filepath.Join(homeDir, ".steam/root"),
		filepath.Join(homeDir, ".steam/debian-installation"),
		// Snap steam install directories
		filepath.Join(snapDir, "steam/common/.local/share/Steam"),
		filepath.Join(snapDir, "steam/common/.steam/steam"),
		filepath.Join(snapDir, "steam/common/.steam/root"),
	}

	// Deduplicate symlinks
	seen := make(map[string]bool)
	var uniquePaths []string

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			continue
		}

		// Resolve symlinks
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			resolved = path
		}

		if !seen[resolved] {
			seen[resolved] = true
			uniquePaths = append(uniquePaths, path)
		}
	}

	return uniquePaths, nil
}
