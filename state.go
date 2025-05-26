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
	"time"
)

// RAIDState represents the state of a RAID configuration.
type RAIDState struct {
	Name      string    `json:"name"`
	RAIDType  string    `json:"raid_type"`
	Drives    []string  `json:"drives"`
	CreatedAt time.Time `json:"created_at"`
}

// SaveRAIDState saves the RAID state to a JSON file and generates a checksum.
// It creates two files: state.json and state.checksum in the specified directory.
func SaveRAIDState(dir string, state RAIDState) error {
	statePath := filepath.Join(dir, "state.json")
	checksumPath := filepath.Join(dir, "state.checksum")

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return err
	}

	sum := sha256.Sum256(data)
	checksum := hex.EncodeToString(sum[:])

	if err := os.WriteFile(checksumPath, []byte(checksum), 0644); err != nil {
		return err
	}

	return nil
}

// LoadRAIDState loads the RAID state from a JSON file and verifies its checksum.
// It expects two files: state.json and state.checksum in the specified directory.
func LoadRAIDState(dir string) (*RAIDState, error) {
	statePath := filepath.Join(dir, "state.json")
	checksumPath := filepath.Join(dir, "state.checksum")

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}

	checksumData, err := os.ReadFile(checksumPath)
	if err != nil {
		return nil, err
	}

	actualChecksum := sha256.Sum256(data)
	if hex.EncodeToString(actualChecksum[:]) != strings.TrimSpace(string(checksumData)) {
		return nil, errors.New("state.json checksum mismatch")
	}

	var state RAIDState

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// ValidateDrives checks if all drives listed in the RAID state exist in the specified base directory.
// It returns an error if any drive is missing.
func ValidateDrives(baseDir string, state *RAIDState) error {
	missing := []string{}

	for _, driveName := range state.Drives {
		path := filepath.Join(baseDir, "drives", driveName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, driveName)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing drives: %v", missing)
	}

	return nil
}
