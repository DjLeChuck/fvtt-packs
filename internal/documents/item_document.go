package documents

import (
	"github.com/djlechuck/fvtt-packs/internal/fvttdb"
)

type ItemDocument struct {
	baseDocument `yaml:",inline"`
	Type         string                  `json:"type" yaml:"type"`
	Img          string                  `json:"img" yaml:"img"`
	System       *System                 `json:"system" yaml:"system"`
	Effects      *[]ActiveEffectDocument `json:"effects" yaml:"effects"`
	Folder       string                  `json:"folder" yaml:"folder"`
	Sort         int                     `json:"sort" yaml:"sort"`
	Ownership    *Ownership              `json:"ownership" yaml:"ownership"`
	Flags        *Flags                  `json:"flags" yaml:"flags"`
	Stats        *DocumentStats          `json:"_stats" yaml:"_stats"`
}

func (d *ItemDocument) HydrateCollections(fvttdb *fvttdb.FvttDb) error {
	return nil
}
