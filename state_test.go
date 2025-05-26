package raid

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadRAIDState_Success(t *testing.T) {
	dir := t.TempDir()
	state := RAIDState{
		RAIDType:  "RAID5",
		Drives:    []string{"sda", "sdb", "sdc"},
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	// Write state.json
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal state: %v", err)
	}
	statePath := filepath.Join(dir, "state.json")
	if err := os.WriteFile(statePath, data, 0644); err != nil {
		t.Fatalf("failed to write state.json: %v", err)
	}

	// Write state.checksum
	sum := sha256.Sum256(data)
	checksum := hex.EncodeToString(sum[:])
	checksumPath := filepath.Join(dir, "state.checksum")
	if err := os.WriteFile(checksumPath, []byte(checksum), 0644); err != nil {
		t.Fatalf("failed to write state.checksum: %v", err)
	}

	loaded, err := LoadRAIDState(dir)
	if err != nil {
		t.Fatalf("LoadRAIDState failed: %v", err)
	}
	if loaded.RAIDType != state.RAIDType {
		t.Errorf("expected RAIDType %q, got %q", state.RAIDType, loaded.RAIDType)
	}
	if len(loaded.Drives) != len(state.Drives) {
		t.Errorf("expected Drives length %d, got %d", len(state.Drives), len(loaded.Drives))
	}
	for i := range state.Drives {
		if loaded.Drives[i] != state.Drives[i] {
			t.Errorf("expected Drives[%d] %q, got %q", i, state.Drives[i], loaded.Drives[i])
		}
	}
	if !loaded.CreatedAt.Equal(state.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", state.CreatedAt, loaded.CreatedAt)
	}
}

func TestLoadRAIDState_ChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	state := RAIDState{
		RAIDType:  "RAID1",
		Drives:    []string{"sda", "sdb"},
		CreatedAt: time.Now().UTC(),
	}
	data, _ := json.MarshalIndent(state, "", "  ")
	statePath := filepath.Join(dir, "state.json")
	os.WriteFile(statePath, data, 0644)
	// Write wrong checksum
	os.WriteFile(filepath.Join(dir, "state.checksum"), []byte("deadbeef"), 0644)

	_, err := LoadRAIDState(dir)
	if err == nil || err.Error() != "state.json checksum mismatch" {
		t.Errorf("expected checksum mismatch error, got %v", err)
	}
}

func TestLoadRAIDState_MissingFiles(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadRAIDState(dir)
	if err == nil {
		t.Error("expected error for missing state.json, got nil")
	}

	// Write only state.json
	state := RAIDState{RAIDType: "RAID0"}
	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(filepath.Join(dir, "state.json"), data, 0644)
	_, err = LoadRAIDState(dir)
	if err == nil {
		t.Error("expected error for missing state.checksum, got nil")
	}
}

func TestLoadRAIDState_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	checksumPath := filepath.Join(dir, "state.checksum")

	invalidJSON := []byte("{invalid json}")
	os.WriteFile(statePath, invalidJSON, 0644)
	sum := sha256.Sum256(invalidJSON)
	os.WriteFile(checksumPath, []byte(hex.EncodeToString(sum[:])), 0644)

	_, err := LoadRAIDState(dir)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
func TestValidateDrives_AllPresent(t *testing.T) {
	dir := t.TempDir()
	drivesDir := filepath.Join(dir, "drives")
	if err := os.Mkdir(drivesDir, 0755); err != nil {
		t.Fatalf("failed to create drives dir: %v", err)
	}
	driveNames := []string{"sda", "sdb", "sdc"}
	for _, d := range driveNames {
		if err := os.WriteFile(filepath.Join(drivesDir, d), []byte("dummy"), 0644); err != nil {
			t.Fatalf("failed to create drive file %s: %v", d, err)
		}
	}
	state := &RAIDState{
		Drives: driveNames,
	}
	if err := ValidateDrives(dir, state); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateDrives_MissingSomeDrives(t *testing.T) {
	dir := t.TempDir()
	drivesDir := filepath.Join(dir, "drives")
	if err := os.Mkdir(drivesDir, 0755); err != nil {
		t.Fatalf("failed to create drives dir: %v", err)
	}
	// Only create sda
	if err := os.WriteFile(filepath.Join(drivesDir, "sda"), []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create drive file sda: %v", err)
	}
	state := &RAIDState{
		Drives: []string{"sda", "sdb", "sdc"},
	}
	err := ValidateDrives(dir, state)
	if err == nil {
		t.Fatal("expected error for missing drives, got nil")
	}
	if got, _ := err.Error(), "missing drives: [sdb sdc]"; !strings.Contains(got, "sdb") || !strings.Contains(got, "sdc") {
		t.Errorf("expected error mentioning missing drives, got %v", got)
	}
}

func TestValidateDrives_AllMissing(t *testing.T) {
	dir := t.TempDir()
	// No drives directory at all
	state := &RAIDState{
		Drives: []string{"sda", "sdb"},
	}
	err := ValidateDrives(dir, state)
	if err == nil {
		t.Fatal("expected error for all missing drives, got nil")
	}
	if got := err.Error(); !strings.Contains(got, "sda") || !strings.Contains(got, "sdb") {
		t.Errorf("expected error mentioning all missing drives, got %v", got)
	}
}

func TestValidateDrives_EmptyDrivesList(t *testing.T) {
	dir := t.TempDir()
	state := &RAIDState{
		Drives: []string{},
	}
	if err := ValidateDrives(dir, state); err != nil {
		t.Errorf("expected no error for empty drives list, got %v", err)
	}
}
