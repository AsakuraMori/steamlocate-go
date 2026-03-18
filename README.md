# steamlocate-go

A Go library for efficiently locating any Steam application on the filesystem, and/or the Steam installation itself.

This is a Go port of the Rust [steamlocate-rs](https://github.com/WilliamVenner/steamlocate-rs) library. Write by AI, just test on Windows.

## Features

- Locate Steam installation directory on Windows, macOS, and Linux
- Find installed Steam apps by App ID
- Iterate over all Steam libraries and apps
- Parse app manifest files (appmanifest_*.acf)
- Read compatibility tool mappings (Proton configurations)
- Parse shortcuts (non-Steam games)

## Installation

```bash
go get github.com/AsakuraMori/steamlocate-go
```
If failed, try to clone this repo.

## Usage

### Locate Steam installation and a specific game

```go
package main

import (
    "fmt"
    "log"
    
    steamlocate "github.com/AsakuraMori/steamlocate-go"
)

func main() {
    steamDir, err := steamlocate.Locate()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Steam installation:", steamDir.Path())
    
    const gmodAppID = 4000
    gmod, library, err := steamDir.FindApp(gmodAppID)
    if err != nil {
        log.Fatal(err)
    }
    
    if gmod != nil {
        fmt.Printf("Found %s in %s\n", gmod.Name, library.Path())
    }
}
```

### Iterate over all libraries and apps

```go
steamDir, err := steamlocate.Locate()
if err != nil {
    log.Fatal(err)
}

libraries, err := steamDir.Libraries()
if err != nil {
    log.Fatal(err)
}

for _, library := range libraries {
    fmt.Println("Library:", library.Path())
    
    apps, err := library.Apps()
    if err != nil {
        continue
    }
    
    for _, app := range apps {
        fmt.Printf("  App %d - %s\n", app.AppID, app.Name)
    }
}
```

## License

MIT License

## Thanks

https://github.com/WilliamVenner/steamlocate-rs
