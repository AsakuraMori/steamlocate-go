package steamlocate

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// VDFNode represents a node in the VDF tree
type VDFNode struct {
	Value    string
	Children map[string]*VDFNode
}

// Get returns a child node by key
func (n *VDFNode) Get(key string) *VDFNode {
	if n == nil || n.Children == nil {
		return nil
	}
	return n.Children[key]
}

// GetString returns the value as a string
func (n *VDFNode) GetString() string {
	if n == nil {
		return ""
	}
	return n.Value
}

// GetInt returns the value as an int
func (n *VDFNode) GetInt() int {
	if n == nil {
		return 0
	}
	i, _ := strconv.Atoi(n.Value)
	return i
}

// GetUint32 returns the value as a uint32
func (n *VDFNode) GetUint32() uint32 {
	if n == nil {
		return 0
	}
	i, _ := strconv.ParseUint(n.Value, 10, 32)
	return uint32(i)
}

// GetUint64 returns the value as a uint64
func (n *VDFNode) GetUint64() uint64 {
	if n == nil {
		return 0
	}
	i, _ := strconv.ParseUint(n.Value, 10, 64)
	return i
}

// GetMap returns the children as a map
func (n *VDFNode) GetMap() map[string]*VDFNode {
	if n == nil {
		return nil
	}
	return n.Children
}

// ParseVDF parses a VDF formatted string
func ParseVDF(data string) (*VDFNode, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanRunes)

	root := &VDFNode{Children: make(map[string]*VDFNode)}
	stack := []*VDFNode{root}
	var currentKey string
	var inKey bool
	var inValue bool
	var escape bool
	var buffer strings.Builder

	for scanner.Scan() {
		char := scanner.Text()

		if escape {
			buffer.WriteString(char)
			escape = false
			continue
		}

		if char == "\\" {
			escape = true
			buffer.WriteString(char)
			continue
		}

		if inValue {
			if char == "\"" {
				inValue = false
				value := buffer.String()
				buffer.Reset()

				// Add key-value pair to current node
				if len(stack) > 0 {
					parent := stack[len(stack)-1]
					if parent.Children == nil {
						parent.Children = make(map[string]*VDFNode)
					}
					parent.Children[currentKey] = &VDFNode{Value: value}
				}
				currentKey = ""
			} else {
				buffer.WriteString(char)
			}
		} else if inKey {
			if char == "\"" {
				inKey = false
				currentKey = buffer.String()
				buffer.Reset()

				// Check if next non-whitespace char is '{' (subtree) or '"' (value)
				// We'll handle this in the main loop by looking ahead
			} else {
				buffer.WriteString(char)
			}
		} else {
			switch char {
			case "\"":
				if currentKey == "" {
					inKey = true
				} else {
					inValue = true
				}
			case "{":
				// Start of a subtree
				if currentKey != "" && len(stack) > 0 {
					parent := stack[len(stack)-1]
					if parent.Children == nil {
						parent.Children = make(map[string]*VDFNode)
					}
					child := &VDFNode{Children: make(map[string]*VDFNode)}
					parent.Children[currentKey] = child
					stack = append(stack, child)
					currentKey = ""
				}
			case "}":
				// End of a subtree
				if len(stack) > 1 {
					stack = stack[:len(stack)-1]
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return root, nil
}

// ParseLibraryFolders parses the libraryfolders.vdf file
func ParseLibraryFolders(data string) ([]string, error) {
	root, err := ParseVDF(data)
	if err != nil {
		return nil, err
	}

	libraryfolders := root.Get("libraryfolders")
	if libraryfolders == nil {
		return nil, fmt.Errorf("missing libraryfolders key")
	}

	var paths []string
	for key, node := range libraryfolders.Children {
		// Keys are numeric indices "0", "1", etc.
		if _, err := strconv.Atoi(key); err == nil {
			if pathNode := node.Get("path"); pathNode != nil {
				paths = append(paths, pathNode.GetString())
			}
		}
	}

	return paths, nil
}
