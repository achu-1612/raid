package raid

import (
	"errors"
	"strings"
)

// RAID5 implements striping with distributed parity
type RAID5 struct {
	name   string
	drives []*Drive
}

func NewRAID5(name string, drives []*Drive) (*RAID5, error) {
	if len(drives) < 3 {
		return nil, errors.New("RAID 5 requires at least three drives")
	}
	return &RAID5{name: name, drives: drives}, nil
}

func (r *RAID5) Write(filename, data string) error {
	n := len(r.drives)
	chunkSize := len(data) / (n - 1)
	if len(data)%(n-1) != 0 {
		chunkSize++
	}
	chunks := make([]string, n)
	for i := 0; i < n; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(data) {
			end = len(data)
		}
		if i < n-1 {
			chunks[i] = data[start:end]
		}
	}
	parityDrive := len(filename) % n
	parity := byte(0)
	for i := 0; i < n-1; i++ {
		for _, b := range []byte(chunks[i]) {
			parity ^= b
		}
	}
	for i := 0; i < n; i++ {
		if i == parityDrive {
			r.drives[i].WriteFile(filename+".parity", string(parity))
		} else {
			r.drives[i].WriteFile(filename, chunks[i])
		}
	}
	return nil
}

func (r *RAID5) Read(filename string) (string, error) {
	var result strings.Builder
	for _, d := range r.drives {
		chunk, err := d.ReadFile(filename)
		if err == nil {
			result.WriteString(chunk)
		}
	}
	return result.String(), nil
}

func (r *RAID5) Name() string { return r.name }

func (r *RAID5) Reconstruct(filename, failedDriveName string) error {
	n := len(r.drives)
	missingIndex := -1
	chunks := make([][]byte, n)
	parity := byte(0)
	parityDrive := len(filename) % n
	for i, d := range r.drives {
		if d.Name == failedDriveName {
			missingIndex = i
			continue
		}
		if i == parityDrive {
			p, _ := d.ReadFile(filename + ".parity")
			if len(p) > 0 {
				parity ^= p[0]
			}
		} else {
			data, _ := d.ReadFile(filename)
			chunks[i] = []byte(data)
			for _, b := range data {
				parity ^= byte(b)
			}
		}
	}
	if missingIndex == -1 {
		return errors.New("failed drive not found")
	}
	recovered := make([]byte, len(chunks[0]))
	for i := range recovered {
		recovered[i] = parity
	}
	d := r.drives[missingIndex]
	d.Recreate()
	return d.WriteFile(filename, string(recovered))
}
