//go:build darwin

package steamlocate

import (
	"os"
	"path/filepath"
)

// locateSteamDir returns the Steam installation directory on macOS
func locateSteamDir() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, newNoHomeError()
	}

	installPath := filepath.Join(homeDir, "Library/Application Support/Steam")
	return []string{installPath}, nil
}
