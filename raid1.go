package raid

import "errors"

// RAID1 implements mirroring
type RAID1 struct {
	name   string
	drives []*Drive
}

func NewRAID1(name string, drives []*Drive) (*RAID1, error) {
	if len(drives) < 2 {
		return nil, errors.New("RAID 1 requires at least two drives")
	}
	return &RAID1{name: name, drives: drives}, nil
}

func (r *RAID1) Type() RAIDType {
	return RAIDType1
}

func (r *RAID1) Write(filename, data string) error {
	for _, d := range r.drives {
		if err := d.WriteFile(filename, data); err != nil {
			return err
		}
	}
	return nil
}

func (r *RAID1) Read(filename string) (string, error) {
	var lastErr error
	for _, d := range r.drives {
		data, err := d.ReadFile(filename)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func (r *RAID1) Name() string { return r.name }

func (r *RAID1) Reconstruct(filename, failedDriveName string) error {
	var refData string
	for _, d := range r.drives {
		if d.Exists(filename) {
			data, err := d.ReadFile(filename)
			if err == nil {
				refData = data
				break
			}
		}
	}
	if refData == "" {
		return errors.New("no source to reconstruct from")
	}
	for _, d := range r.drives {
		if d.name == failedDriveName {
			d.Recreate()
			return d.WriteFile(filename, refData)
		}
	}
	return errors.New("failed drive not found")
}
