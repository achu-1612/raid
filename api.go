package raid

// RAID interface defines the methods for RAID implementations
type RAID interface {
	Write(filename, data string) error
	Read(filename string) (string, error)
	Name() string
	Reconstruct(filename, failedDriveName string) error
	Type() RAIDType
}

// Storage interface defines the methods for storage operations
// that RAID implementations will use to interact with drives
type Storage interface {
	Name() string
	WriteFile(filename, data string) error
	ReadFile(filename string) (string, error)
	Exists(filename string) bool
	Recreate() error
}

// rType represents the type of RAID
type RAIDType string

// Supported RAID types
const (
	RAIDType0  RAIDType = "RAID0"
	RAIDType1  RAIDType = "RAID1"
	RAIDType5  RAIDType = "RAID5"
	RAIDType10 RAIDType = "RAID10"
)

func (r RAIDType) IsValid() bool {
	switch r {
	case RAIDType0, RAIDType1, RAIDType5, RAIDType10:
		return true
	default:
		return false
	}
}
