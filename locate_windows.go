//go:build windows

package steamlocate

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procRegOpenKeyExW        = advapi32.NewProc("RegOpenKeyExW")
	procRegQueryValueExW     = advapi32.NewProc("RegQueryValueExW")
	procRegCloseKey          = advapi32.NewProc("RegCloseKey")
)

const (
	hkeyLocalMachine = 0x80000002
	keyQueryValue    = 0x0001
	keyWow6464Key    = 0x0100
	keyWow6432Key    = 0x0200
)

type registryKey syscall.Handle

// locateSteamDir returns the Steam installation directory on Windows
func locateSteamDir() ([]string, error) {
	// Try 32-bit registry first (for 32-bit Steam on 64-bit Windows), then 64-bit
	installPath, err := getInstallPathFromRegistry(keyWow6432Key)
	if err != nil {
		installPath, err = getInstallPathFromRegistry(keyWow6464Key)
		if err != nil {
			// Try without WOW64 flags
			installPath, err = getInstallPathFromRegistry(0)
			if err != nil {
				return nil, newWinRegError(err)
			}
		}
	}

	return []string{installPath}, nil
}

func getInstallPathFromRegistry(wow64 uint32) (string, error) {
	// Try standard path first
	path, err := getRegStringValue(hkeyLocalMachine, `SOFTWARE\Valve\Steam`, "InstallPath", wow64)
	if err != nil {
		// Try Wow6432Node for 32-bit Steam on 64-bit Windows
		path, err = getRegStringValue(hkeyLocalMachine, `SOFTWARE\Wow6432Node\Valve\Steam`, "InstallPath", wow64)
		if err != nil {
			return "", err
		}
	}
	return filepath.Clean(path), nil
}

func getRegStringValue(root uint32, path string, valueName string, wow64 uint32) (string, error) {
	var keyHandle registryKey
	pathPtr, _ := syscall.UTF16PtrFromString(path)
	
	access := keyQueryValue | wow64
	
	ret, _, err := procRegOpenKeyExW.Call(
		uintptr(root),
		uintptr(unsafe.Pointer(pathPtr)),
		0,
		uintptr(access),
		uintptr(unsafe.Pointer(&keyHandle)),
	)
	
	if ret != 0 {
		return "", err
	}
	defer procRegCloseKey.Call(uintptr(keyHandle))

	valuePtr, _ := syscall.UTF16PtrFromString(valueName)
	
	var bufSize uint32
	var valueType uint32
	
	// First call to get size
	ret, _, _ = procRegQueryValueExW.Call(
		uintptr(keyHandle),
		uintptr(unsafe.Pointer(valuePtr)),
		0,
		uintptr(unsafe.Pointer(&valueType)),
		0,
		uintptr(unsafe.Pointer(&bufSize)),
	)
	
	if ret != 0 {
		return "", syscall.Errno(ret)
	}
	
	if valueType != syscall.REG_SZ && valueType != syscall.REG_EXPAND_SZ {
		return "", syscall.EINVAL
	}
	
	// Allocate buffer
	buf := make([]uint16, bufSize/2)
	
	// Second call to get data
	ret, _, _ = procRegQueryValueExW.Call(
		uintptr(keyHandle),
		uintptr(unsafe.Pointer(valuePtr)),
		0,
		uintptr(unsafe.Pointer(&valueType)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufSize)),
	)
	
	if ret != 0 {
		return "", syscall.Errno(ret)
	}
	
	return syscall.UTF16ToString(buf), nil
}
