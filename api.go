package raid

// RAID interface defines the methods for RAID implementations
type RAID interface {
	Write(filename, data string) error
	Read(filename string) (string, error)
	Name() string
	Reconstruct(filename, failedDriveName string) error
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
type rType string

// Supported RAID types
const (
	RAIDType0  rType = "raid0"
	RAIDType1  rType = "raid1"
	RAIDType5  rType = "raid5"
	RAIDType10 rType = "raid10"
)

func (r rType) IsValid() bool {
	switch r {
	case RAIDType0, RAIDType1, RAIDType5, RAIDType10:
		return true
	default:
		return false
	}
}
