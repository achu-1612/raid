package raid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const baseDir = "./vfs"

// Drive represents a virtual file system
type Drive struct {
	name   string
	path   string
	lock   sync.RWMutex
	failed bool
}

func createDrive(name string) (*Drive, error) {
	drivePath := filepath.Join(baseDir, "drives", name)
	if err := os.MkdirAll(drivePath, 0755); err != nil {
		return nil, err
	}
	return &Drive{name: name, path: drivePath}, nil
}

func createDrivesForRAID(name string, raidType RAIDType) ([]*Drive, error) {
	var count int
	switch raidType {
	case RAIDType0:
		count = 2

	case RAIDType1:
		count = 2

	case RAIDType5:
		count = 3

	case RAIDType10:
		count = 4

	default:
		return nil, fmt.Errorf("unsupported raid type: %s", raidType)
	}

	var drives []*Drive
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%s%d", name, i)

		drive, err := createDrive(name)
		if err != nil {
			return nil, fmt.Errorf("failed to create drive %s: %w", name, err)
		}

		drives = append(drives, drive)
	}

	return drives, nil
}

func (d *Drive) WriteFile(filename, data string) error {
	if d.failed {
		return errors.New("drive failed")
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	return os.WriteFile(filepath.Join(d.path, filename), []byte(data), 0644)
}

func (d *Drive) ReadFile(filename string) (string, error) {
	if d.failed {
		return "", errors.New("drive failed")
	}
	d.lock.RLock()
	defer d.lock.RUnlock()
	data, err := os.ReadFile(filepath.Join(d.path, filename))
	return string(data), err
}

func (d *Drive) Exists(filename string) bool {
	_, err := os.Stat(filepath.Join(d.path, filename))
	return err == nil
}

func (d *Drive) Recreate() error {
	return os.MkdirAll(d.path, 0755)
}

func (d *Drive) Name() string {
	return d.name
}

func drivesToString(drives []*Drive) []string {
	var names []string
	for _, drive := range drives {
		names = append(names, drive.Name())
	}
	return names
}
