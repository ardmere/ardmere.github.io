package sitegen

import (
	"os"

	"gopkg.in/yaml.v3"
)

const defaultSnapshotNotesPath = "config/site/snapshot-notes.yaml"

type snapshotNotesFile struct {
	Overrides map[string]string `yaml:"overrides"`
	Footers   []string          `yaml:"footers"`
}

func loadSnapshotNotes() (*snapshotNotesFile, error) {
	path := defaultSnapshotNotesPath
	if _, err := os.Stat(path); err != nil {
		return nil, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var notes snapshotNotesFile
	if err := yaml.Unmarshal(raw, &notes); err != nil {
		return nil, err
	}
	if notes.Overrides == nil {
		notes.Overrides = map[string]string{}
	}
	return &notes, nil
}

func snapshotFooters() []string {
	notes, err := loadSnapshotNotes()
	if err != nil || notes == nil {
		return nil
	}
	return notes.Footers
}
