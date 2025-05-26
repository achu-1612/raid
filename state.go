package raid

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RAIDState represents the state of a RAID configuration.
type RAIDState struct {
	Name     string   `json:"name"`
	RAIDType string   `json:"raid_type"`
	Drives   []string `json:"drives"`
}

const (
	// StateFileName is the name of the file where RAID state is stored.
	StateFileName = ".state.json"
	// ChecksumFileName is the name of the file where the checksum of the RAID state is stored.
	ChecksumFileName = ".state.checksum"
)

// InitializeState creates the necessary files for RAID state management.
func InitializeState(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	statePath := filepath.Join(dir, StateFileName)
	checksumPath := filepath.Join(dir, ChecksumFileName)

	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		return fmt.Errorf("state file already exists at %s", statePath)
	}

	if _, err := os.Stat(checksumPath); !os.IsNotExist(err) {
		return fmt.Errorf("checksum file already exists at %s", checksumPath)
	}

	if err := DumpState(dir, map[string]RAIDState{}); err != nil {
		return fmt.Errorf("store empty raid state: %w", err)
	}

	return nil
}

// SaveRAIDState saves the RAID state to a JSON file and generates a checksum.
// It updates two files: .state.json and .state.checksum in the specified directory.
func SaveRAIDState(dir string, state RAIDState) error {
	currentState, err := LoadRAIDState(dir)
	if err != nil {
		return fmt.Errorf("SaveRAIDState: %w", err)
	}

	if _, ok := currentState[state.Name]; ok {
		return fmt.Errorf("RAID state with name %s already exists", state.Name)
	}

	currentState[state.Name] = state

	if err := DumpState(dir, currentState); err != nil {
		return fmt.Errorf("DumpState: %v", err)
	}

	fmt.Printf("RAID state for %s saved successfully.\n", state.Name)

	return nil
}

// LoadRAIDState loads the RAID state from a JSON file and verifies its checksum.
// It expects two files: .state.json and .state.checksum in the specified directory.
func LoadRAIDState(dir string) (map[string]RAIDState, error) {
	statePath := filepath.Join(dir, StateFileName)
	checksumPath := filepath.Join(dir, ChecksumFileName)

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	checksumData, err := os.ReadFile(checksumPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read checksum file: %w", err)
	}

	actualChecksum := sha256.Sum256(data)
	if hex.EncodeToString(actualChecksum[:]) != strings.TrimSpace(string(checksumData)) {
		return nil, errors.New("checksum mismatch")
	}

	var state map[string]RAIDState

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RAID state: %w", err)
	}

	return state, nil
}

// ValidateDrives checks if all drives listed in the RAID state exist in the specified base directory.
// It returns an error if any drive is missing.
func ValidateDrives(dir string, state *RAIDState) error {
	missing := []string{}

	for _, driveName := range state.Drives {
		path := filepath.Join(dir, "drives", driveName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, driveName)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing drives: %v", missing)
	}

	return nil
}

// DumpState saves the RAID state to a JSON file and generates a checksum.
func DumpState(dir string, state map[string]RAIDState) error {
	statePath := filepath.Join(dir, StateFileName)
	checksumPath := filepath.Join(dir, ChecksumFileName)

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RAID state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write RAID state file: %w", err)
	}

	sum := sha256.Sum256(data)
	checksum := hex.EncodeToString(sum[:])

	if err := os.WriteFile(checksumPath, []byte(checksum), 0644); err != nil {
		return fmt.Errorf("failed to write checksum file: %w", err)
	}

	return nil
}

func RAIDExists(dir, name string) (bool, error) {
	state, err := LoadRAIDState(dir)
	if err != nil {
		return false, fmt.Errorf("failed to load RAID state: %w", err)
	}

	_, exists := state[name]
	return exists, nil
}
