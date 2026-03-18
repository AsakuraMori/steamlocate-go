package steamlocate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Universe represents the Steam universe
type Universe int

const (
	UniverseInvalid Universe = iota
	UniversePublic
	UniverseBeta
	UniverseInternal
	UniverseDev
)

func (u Universe) String() string {
	switch u {
	case UniverseInvalid:
		return "Invalid"
	case UniversePublic:
		return "Public"
	case UniverseBeta:
		return "Beta"
	case UniverseInternal:
		return "Internal"
	case UniverseDev:
		return "Dev"
	default:
		return fmt.Sprintf("Unknown(%d)", u)
	}
}

// StateFlag represents a state flag for an app
type StateFlag int

const (
	StateFlagInvalid StateFlag = iota
	StateFlagUninstalled
	StateFlagUpdateRequired
	StateFlagFullyInstalled
	StateFlagEncrypted
	StateFlagLocked
	StateFlagFilesMissing
	StateFlagAppRunning
	StateFlagFilesCorrupt
	StateFlagUpdateRunning
	StateFlagUpdatePaused
	StateFlagUpdateStarted
	StateFlagUninstalling
	StateFlagBackupRunning
	StateFlagReconfiguring
	StateFlagValidating
	StateFlagAddingFiles
	StateFlagPreallocating
	StateFlagDownloading
	StateFlagStaging
	StateFlagCommitting
	StateFlagUpdateStopping
)

func (sf StateFlag) String() string {
	names := map[StateFlag]string{
		StateFlagInvalid:          "Invalid",
		StateFlagUninstalled:      "Uninstalled",
		StateFlagUpdateRequired:   "UpdateRequired",
		StateFlagFullyInstalled:   "FullyInstalled",
		StateFlagEncrypted:        "Encrypted",
		StateFlagLocked:           "Locked",
		StateFlagFilesMissing:     "FilesMissing",
		StateFlagAppRunning:       "AppRunning",
		StateFlagFilesCorrupt:     "FilesCorrupt",
		StateFlagUpdateRunning:    "UpdateRunning",
		StateFlagUpdatePaused:     "UpdatePaused",
		StateFlagUpdateStarted:    "UpdateStarted",
		StateFlagUninstalling:     "Uninstalling",
		StateFlagBackupRunning:    "BackupRunning",
		StateFlagReconfiguring:    "Reconfiguring",
		StateFlagValidating:       "Validating",
		StateFlagAddingFiles:      "AddingFiles",
		StateFlagPreallocating:    "Preallocating",
		StateFlagDownloading:      "Downloading",
		StateFlagStaging:          "Staging",
		StateFlagCommitting:       "Committing",
		StateFlagUpdateStopping:   "UpdateStopping",
	}
	if name, ok := names[sf]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", sf)
}

// StateFlags holds the state flags value and can iterate over individual flags
type StateFlags uint64

// Flags returns an iterator (slice) of individual flags
func (sf StateFlags) Flags() []StateFlag {
	if sf == 0 {
		return []StateFlag{StateFlagInvalid}
	}

	var flags []StateFlag
	for i := 0; i < 64; i++ {
		if sf&(1<<i) != 0 {
			flag := stateFlagFromBitOffset(i)
			flags = append(flags, flag)
		}
	}
	return flags
}

func stateFlagFromBitOffset(offset int) StateFlag {
	switch offset {
	case 0:
		return StateFlagUninstalled
	case 1:
		return StateFlagUpdateRequired
	case 2:
		return StateFlagFullyInstalled
	case 3:
		return StateFlagEncrypted
	case 4:
		return StateFlagLocked
	case 5:
		return StateFlagFilesMissing
	case 6:
		return StateFlagAppRunning
	case 7:
		return StateFlagFilesCorrupt
	case 8:
		return StateFlagUpdateRunning
	case 9:
		return StateFlagUpdatePaused
	case 10:
		return StateFlagUpdateStarted
	case 11:
		return StateFlagUninstalling
	case 12:
		return StateFlagBackupRunning
	case 16:
		return StateFlagReconfiguring
	case 17:
		return StateFlagValidating
	case 18:
		return StateFlagAddingFiles
	case 19:
		return StateFlagPreallocating
	case 20:
		return StateFlagDownloading
	case 21:
		return StateFlagStaging
	case 22:
		return StateFlagCommitting
	case 23:
		return StateFlagUpdateStopping
	default:
		return StateFlagInvalid
	}
}

// AutoUpdateBehavior represents auto update behavior
type AutoUpdateBehavior int

const (
	AutoUpdateBehaviorKeepUpToDate AutoUpdateBehavior = iota
	AutoUpdateBehaviorOnlyUpdateOnLaunch
	AutoUpdateBehaviorUpdateWithHighPriority
)

func (aub AutoUpdateBehavior) String() string {
	switch aub {
	case AutoUpdateBehaviorKeepUpToDate:
		return "KeepUpToDate"
	case AutoUpdateBehaviorOnlyUpdateOnLaunch:
		return "OnlyUpdateOnLaunch"
	case AutoUpdateBehaviorUpdateWithHighPriority:
		return "UpdateWithHighPriority"
	default:
		return fmt.Sprintf("Unknown(%d)", aub)
	}
}

// AllowOtherDownloadsWhileRunning represents download behavior
type AllowOtherDownloadsWhileRunning int

const (
	AllowOtherDownloadsUseGlobalSetting AllowOtherDownloadsWhileRunning = iota
	AllowOtherDownloadsAllow
	AllowOtherDownloadsNever
)

func (aod AllowOtherDownloadsWhileRunning) String() string {
	switch aod {
	case AllowOtherDownloadsUseGlobalSetting:
		return "UseGlobalSetting"
	case AllowOtherDownloadsAllow:
		return "Allow"
	case AllowOtherDownloadsNever:
		return "Never"
	default:
		return fmt.Sprintf("Unknown(%d)", aod)
	}
}

// Depot represents a depot entry
type Depot struct {
	Manifest  uint64
	Size      uint64
	DLCAppID  *uint64
}

// App represents a Steam application manifest
type App struct {
	AppID                          uint32
	InstallDir                     string
	Name                           string
	LastUser                       *uint64
	Universe                       *Universe
	LauncherPath                   string
	StateFlags                     *StateFlags
	LastUpdated                    *time.Time
	UpdateResult                   *uint64
	SizeOnDisk                     *uint64
	BuildID                        *uint64
	BytesToDownload                *uint64
	BytesDownloaded                *uint64
	BytesToStage                   *uint64
	BytesStaged                    *uint64
	StagingSize                    *uint64
	TargetBuildID                  *uint64
	AutoUpdateBehavior             *AutoUpdateBehavior
	AllowOtherDownloadsWhileRunning *AllowOtherDownloadsWhileRunning
	ScheduledAutoUpdate            *time.Time
	FullValidateBeforeNextUpdate   *bool
	FullValidateAfterNextUpdate    *bool
	InstalledDepots                map[uint64]*Depot
	StagedDepots                   map[uint64]*Depot
	UserConfig                     map[string]string
	MountedConfig                  map[string]string
	InstallScripts                 map[uint64]string
	SharedDepots                   map[uint64]uint64
}

// ParseAppManifest parses an appmanifest_*.acf file
func ParseAppManifest(path string) (*App, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, newIOError(path, err)
	}

	root, err := ParseVDF(string(data))
	if err != nil {
		return nil, newParseError(ParseErrorKindApp, path, err)
	}

	appState := root.Get("AppState")
	if appState == nil {
		return nil, newUnexpectedStructureError(ParseErrorKindApp, path)
	}

	app := &App{
		InstalledDepots: make(map[uint64]*Depot),
		StagedDepots:    make(map[uint64]*Depot),
		UserConfig:      make(map[string]string),
		MountedConfig:   make(map[string]string),
		InstallScripts:  make(map[uint64]string),
		SharedDepots:    make(map[uint64]uint64),
	}

	// Parse basic fields
	if v := appState.Get("appid"); v != nil {
		app.AppID = v.GetUint32()
	}
	if v := appState.Get("appid"); v == nil || v.GetString() == "" {
		// Try "appID" variant
		if v := appState.Get("appID"); v != nil {
			app.AppID = v.GetUint32()
		}
	}

	if v := appState.Get("installdir"); v != nil {
		app.InstallDir = v.GetString()
	}
	if v := appState.Get("name"); v != nil {
		app.Name = v.GetString()
	}
	if v := appState.Get("name"); v == nil || v.GetString() == "" {
		// Try "Name" variant
		if v := appState.Get("Name"); v != nil {
			app.Name = v.GetString()
		}
	}

	if v := appState.Get("LastOwner"); v != nil {
		u := v.GetUint64()
		app.LastUser = &u
	}
	if v := appState.Get("Universe"); v != nil {
		u := Universe(v.GetInt())
		app.Universe = &u
	}
	if v := appState.Get("LauncherPath"); v != nil {
		app.LauncherPath = v.GetString()
	}

	if v := appState.Get("StateFlags"); v != nil {
		sf := StateFlags(v.GetUint64())
		app.StateFlags = &sf
	}

	if v := appState.Get("LastUpdated"); v != nil {
		t := parseTimestamp(v.GetString())
		if t != nil {
			app.LastUpdated = t
		}
	}
	if v := appState.Get("lastupdated"); v != nil {
		t := parseTimestamp(v.GetString())
		if t != nil {
			app.LastUpdated = t
		}
	}

	if v := appState.Get("UpdateResult"); v != nil {
		u := v.GetUint64()
		app.UpdateResult = &u
	}
	if v := appState.Get("SizeOnDisk"); v != nil {
		u := v.GetUint64()
		app.SizeOnDisk = &u
	}
	if v := appState.Get("buildid"); v != nil {
		u := v.GetUint64()
		app.BuildID = &u
	}
	if v := appState.Get("BytesToDownload"); v != nil {
		u := v.GetUint64()
		app.BytesToDownload = &u
	}
	if v := appState.Get("BytesDownloaded"); v != nil {
		u := v.GetUint64()
		app.BytesDownloaded = &u
	}
	if v := appState.Get("BytesToStage"); v != nil {
		u := v.GetUint64()
		app.BytesToStage = &u
	}
	if v := appState.Get("BytesStaged"); v != nil {
		u := v.GetUint64()
		app.BytesStaged = &u
	}
	if v := appState.Get("StagingSize"); v != nil {
		u := v.GetUint64()
		app.StagingSize = &u
	}
	if v := appState.Get("TargetBuildID"); v != nil {
		u := v.GetUint64()
		app.TargetBuildID = &u
	}

	if v := appState.Get("AutoUpdateBehavior"); v != nil {
		a := AutoUpdateBehavior(v.GetInt())
		app.AutoUpdateBehavior = &a
	}
	if v := appState.Get("AllowOtherDownloadsWhileRunning"); v != nil {
		a := AllowOtherDownloadsWhileRunning(v.GetInt())
		app.AllowOtherDownloadsWhileRunning = &a
	}

	// Parse ScheduledAutoUpdate
	if v := appState.Get("ScheduledAutoUpdate"); v != nil {
		t := parseTimestamp(v.GetString())
		if t != nil {
			app.ScheduledAutoUpdate = t
		}
	}

	// Parse depots
	if depotsNode := appState.Get("InstalledDepots"); depotsNode != nil {
		for depotIDStr, depotNode := range depotsNode.Children {
			depotID, _ := strconv.ParseUint(depotIDStr, 10, 64)
			depot := &Depot{}
			if v := depotNode.Get("manifest"); v != nil {
				depot.Manifest = v.GetUint64()
			}
			if v := depotNode.Get("size"); v != nil {
				depot.Size = v.GetUint64()
			}
			if v := depotNode.Get("dlcappid"); v != nil {
				dlc := v.GetUint64()
				depot.DLCAppID = &dlc
			}
			app.InstalledDepots[depotID] = depot
		}
	}

	// Parse user config
	if userConfigNode := appState.Get("UserConfig"); userConfigNode != nil {
		for k, v := range userConfigNode.Children {
			app.UserConfig[k] = v.GetString()
		}
	}

	// Parse mounted config
	if mountedConfigNode := appState.Get("MountedConfig"); mountedConfigNode != nil {
		for k, v := range mountedConfigNode.Children {
			app.MountedConfig[k] = v.GetString()
		}
	}

	return app, nil
}

func parseTimestamp(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return nil
	}
	secs, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	t := time.Unix(secs, 0)
	return &t
}

// ResolveAppDir resolves the installation directory for an app within a library
func ResolveAppDir(libraryPath string, app *App) string {
	return filepath.Join(libraryPath, "steamapps", "common", app.InstallDir)
}
