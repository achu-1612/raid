package raid

import (
	"errors"
	"strings"
)

// RAID0 implements striping across drives
type RAID0 struct {
	name   string
	drives []*Drive
}

func NewRAID0(name string, drives []*Drive) (*RAID0, error) {
	if len(drives) < 2 {
		return nil, errors.New("RAID 0 requires at least two drives")
	}
	return &RAID0{name: name, drives: drives}, nil
}

func (r *RAID0) Type() RAIDType {
	return RAIDType0
}

func (r *RAID0) Write(filename, data string) error {
	chunkSize := len(data) / len(r.drives)
	for i, d := range r.drives {
		start := i * chunkSize
		end := start + chunkSize
		if i == len(r.drives)-1 {
			end = len(data)
		}
		if err := d.WriteFile(filename, data[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (r *RAID0) Read(filename string) (string, error) {
	var result strings.Builder
	for _, d := range r.drives {
		chunk, err := d.ReadFile(filename)
		if err != nil {
			return "", err
		}
		result.WriteString(chunk)
	}
	return result.String(), nil
}

func (r *RAID0) Name() string { return r.name }

func (r *RAID0) Reconstruct(filename, failedDriveName string) error {
	return errors.New("RAID 0 does not support reconstruction")
}
