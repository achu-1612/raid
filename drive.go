package raid

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

const baseDir = "./vfs"

// Drive represents a virtual file system
type Drive struct {
	Name   string
	Path   string
	lock   sync.RWMutex
	Failed bool
}

func createDrive(name string) (*Drive, error) {
	drivePath := filepath.Join(baseDir, "drives", name)
	if err := os.MkdirAll(drivePath, 0755); err != nil {
		return nil, err
	}
	return &Drive{Name: name, Path: drivePath}, nil
}

func createDrivesForRAID(raidType rType) ([]*Drive, error) {
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

	// Generate a random starting letter from 'a' to 'v' to allow room
	baseRune := rune('a' + rand.Intn(22)) // avoid going past 'z'

	var drives []*Drive

	for i := 0; i < count; i++ {
		driveLetter := rune(int(baseRune) + i)
		if driveLetter > 'z' {
			return nil, errors.New("ran out of drive letters")
		}

		name := fmt.Sprintf("sd%c", driveLetter)

		drive, err := createDrive(name)
		if err != nil {
			return nil, fmt.Errorf("failed to create drive %s: %w", name, err)
		}

		drives = append(drives, drive)
	}

	return drives, nil
}

func (d *Drive) WriteFile(filename, data string) error {
	if d.Failed {
		return errors.New("drive failed")
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	return os.WriteFile(filepath.Join(d.Path, filename), []byte(data), 0644)
}

func (d *Drive) ReadFile(filename string) (string, error) {
	if d.Failed {
		return "", errors.New("drive failed")
	}
	d.lock.RLock()
	defer d.lock.RUnlock()
	data, err := os.ReadFile(filepath.Join(d.Path, filename))
	return string(data), err
}

func (d *Drive) Exists(filename string) bool {
	_, err := os.Stat(filepath.Join(d.Path, filename))
	return err == nil
}

func (d *Drive) Recreate() error {
	return os.MkdirAll(d.Path, 0755)
}
