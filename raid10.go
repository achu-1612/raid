package raid

import (
	"errors"
	"strings"
)

// RAID10 combines mirroring and striping
type RAID10 struct {
	name  string
	pairs [][2]*Drive
}

func NewRAID10(name string, drives []*Drive) (*RAID10, error) {
	if len(drives)%2 != 0 || len(drives) < 4 {
		return nil, errors.New("RAID 10 requires even number of at least 4 drives")
	}
	pairs := make([][2]*Drive, 0)
	for i := 0; i < len(drives); i += 2 {
		pairs = append(pairs, [2]*Drive{drives[i], drives[i+1]})
	}
	return &RAID10{name: name, pairs: pairs}, nil
}

func (r *RAID10) Write(filename, data string) error {
	chunkSize := len(data) / len(r.pairs)
	for i, pair := range r.pairs {
		start := i * chunkSize
		end := start + chunkSize
		if i == len(r.pairs)-1 {
			end = len(data)
		}
		chunk := data[start:end]
		for _, d := range pair {
			if err := d.WriteFile(filename, chunk); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *RAID10) Read(filename string) (string, error) {
	var result strings.Builder
	for _, pair := range r.pairs {
		data, err := pair[0].ReadFile(filename)
		if err != nil {
			data, err = pair[1].ReadFile(filename)
		}
		if err != nil {
			return "", err
		}
		result.WriteString(data)
	}
	return result.String(), nil
}

func (r *RAID10) Name() string { return r.name }

func (r *RAID10) Reconstruct(filename, failedDriveName string) error {
	for _, pair := range r.pairs {
		var ref *Drive
		var target *Drive
		if pair[0].Name == failedDriveName {
			ref, target = pair[1], pair[0]
		} else if pair[1].Name == failedDriveName {
			ref, target = pair[0], pair[1]
		} else {
			continue
		}
		data, err := ref.ReadFile(filename)
		if err != nil {
			return err
		}
		target.Recreate()
		return target.WriteFile(filename, data)
	}
	return errors.New("failed drive not found")
}
