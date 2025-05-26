package raid

import (
	"fmt"
)

func New(raidType rType, name string) (RAID, error) {

	if !raidType.IsValid() {
		return nil, fmt.Errorf("invalid RAID type: %v", raidType)
	}

	drives, err := createDrivesForRAID(raidType)
	if err != nil {
		return nil, fmt.Errorf("failed to create drives for RAID: %w", err)
	}

	switch raidType {
	case RAIDType0:
		return NewRAID0(name, drives)

	case RAIDType1:
		return NewRAID1(name, drives)

	case RAIDType5:
		return NewRAID5(name, drives)

	case RAIDType10:
		return NewRAID10(name, drives)

	default:
		return nil, fmt.Errorf("unsupported RAID type: %v", raidType)
	}
}
