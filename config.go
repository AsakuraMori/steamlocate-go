package steamlocate

import (
	"fmt"
	"os"
	"path/filepath"
)

// CompatTool represents a compatibility tool entry (e.g., Proton configuration)
type CompatTool struct {
	Name     string
	Config   string
	Priority uint64
}

// CompatToolMapping represents the compatibility tool mapping from config.vdf
type CompatToolMapping map[uint32]*CompatTool

// ParseCompatToolMapping parses the config.vdf file for compatibility tool mappings
func ParseCompatToolMapping(steamPath string) (CompatToolMapping, error) {
	configPath := filepath.Join(steamPath, "config", "config.vdf")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty mapping if config doesn't exist
			return make(CompatToolMapping), nil
		}
		return nil, newIOError(configPath, err)
	}

	root, err := ParseVDF(string(data))
	if err != nil {
		return nil, newParseError(ParseErrorKindConfig, configPath, err)
	}

	// Navigate: Software -> Valve -> Steam -> CompatToolMapping
	software := root.Get("Software")
	if software == nil {
		software = root.Get("software")
	}
	if software == nil {
		return make(CompatToolMapping), nil
	}

	valve := software.Get("Valve")
	if valve == nil {
		valve = software.Get("valve")
	}
	if valve == nil {
		return make(CompatToolMapping), nil
	}

	steam := valve.Get("Steam")
	if steam == nil {
		steam = valve.Get("steam")
	}
	if steam == nil {
		return make(CompatToolMapping), nil
	}

	mappingNode := steam.Get("CompatToolMapping")
	if mappingNode == nil {
		return make(CompatToolMapping), nil
	}

	result := make(CompatToolMapping)
	for appIDStr, toolNode := range mappingNode.Children {
		var appID uint32
		if _, err := fmt.Sscanf(appIDStr, "%d", &appID); err != nil {
			continue
		}

		tool := &CompatTool{}
		if nameNode := toolNode.Get("name"); nameNode != nil {
			tool.Name = nameNode.GetString()
		}
		if configNode := toolNode.Get("config"); configNode != nil {
			tool.Config = configNode.GetString()
		}
		if priorityNode := toolNode.Get("priority"); priorityNode != nil {
			tool.Priority = priorityNode.GetUint64()
		}

		result[appID] = tool
	}

	return result, nil
}
