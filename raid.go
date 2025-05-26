package raid

import (
	"fmt"
)

func New(dir string, raidType RAIDType, name string) (RAID, error) {
	if !raidType.IsValid() {
		return nil, fmt.Errorf("invalid RAID type: %v", raidType)
	}

	status, err := RAIDExists(dir, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if RAID exists: %w", err)
	}

	if status {
		return nil, fmt.Errorf("RAID with name %s already exists", name)
	}

	drives, err := createDrivesForRAID(name, raidType)
	if err != nil {
		return nil, fmt.Errorf("failed to create drives for RAID: %w", err)
	}

	var r RAID

	switch raidType {
	case RAIDType0:
		r, err = NewRAID0(name, drives)

	case RAIDType1:
		r, err = NewRAID1(name, drives)

	case RAIDType5:
		r, err = NewRAID5(name, drives)

	case RAIDType10:
		r, err = NewRAID10(name, drives)

	default:
		return nil, fmt.Errorf("unsupported RAID type: %v", raidType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create RAID %s: %w", raidType, err)
	}

	if err := SaveRAIDState(dir, RAIDState{
		Name:     name,
		RAIDType: string(raidType),
		Drives:   drivesToString(drives),
	}); err != nil {
		return nil, fmt.Errorf("failed to save RAID state: %w", err)
	}

	return r, nil
}
