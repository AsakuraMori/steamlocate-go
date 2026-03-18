// Package steamlocate provides functionality for locating Steam installation directories
// and installed Steam applications on Windows, macOS, and Linux.
//
// This is a Go port of the Rust steamlocate library (https://github.com/WilliamVenner/steamlocate-rs).
//
// # Main Types
//
// The primary entry point is the SteamDir type, which represents the Steam installation
// directory:
//
//	steamDir, err := steamlocate.Locate()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Finding Apps
//
// You can find a specific app by its Steam App ID:
//
//	app, library, err := steamDir.FindApp(4000) // Garry's Mod
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if app != nil {
//	    fmt.Printf("Found %s at %s\n", app.Name, library.ResolveAppDir(app))
//	}
//
// # Iterating Libraries
//
// You can iterate over all Steam libraries and their apps:
//
//	libraries, err := steamDir.Libraries()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, library := range libraries {
//	    apps, err := library.Apps()
//	    if err != nil {
//	        continue
//	    }
//	    for _, app := range apps {
//	        fmt.Printf("%d: %s\n", app.AppID, app.Name)
//	    }
//	}
//
// # Compatibility Tools
//
// You can access compatibility tool mappings (e.g., Proton configurations):
//
//	mapping, err := steamDir.CompatToolMapping()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for appID, tool := range mapping {
//	    fmt.Printf("App %d uses %s\n", appID, tool.Name)
//	}
//
// # Non-Steam Shortcuts
//
// You can also access non-Steam game shortcuts:
//
//	shortcuts, err := steamDir.Shortcuts()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, shortcut := range shortcuts {
//	    fmt.Printf("%s: %s\n", shortcut.AppName, shortcut.Executable)
//	}
//
// # Error Handling
//
// The package uses custom error types that can be checked using helper functions:
//
//	err := someOperation()
//	if steamlocate.IsLocateError(err) {
//	    // Failed to locate Steam installation
//	}
//	if steamlocate.IsParseError(err) {
//	    // Failed to parse a VDF file
//	}
package steamlocate
