package documents

import "github.com/djlechuck/fvtt-packs/internal/fvttdb"

type FolderDocument struct {
	baseDocument `yaml:",inline"`
	Type         string         `json:"type" yaml:"type"`
	Description  string         `json:"description" yaml:"description"`
	Folder       string         `json:"folder" yaml:"folder"`
	Sorting      string         `json:"sorting" yaml:"sorting"`
	Sort         int            `json:"sort" yaml:"sort"`
	Color        string         `json:"color" yaml:"color"`
	Flags        *Flags         `json:"flags" yaml:"flags"`
	Stats        *DocumentStats `json:"_stats" yaml:"_stats"`
}

func (d *FolderDocument) HydrateCollections(fvttdb *fvttdb.FvttDb) error {
	return nil
}
