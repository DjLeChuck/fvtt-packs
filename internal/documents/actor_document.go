package documents

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type ActorDocument struct {
	baseDocument   `yaml:",inline"`
	Img            string                  `json:"img" yaml:"img"`
	Type           string                  `json:"type" yaml:"type"`
	System         *System                 `json:"system" yaml:"system"`
	PrototypeToken *PrototypeTokenDocument `json:"prototypeToken" yaml:"prototypeToken"`
	Items          *[]ItemDocument         `json:"-" yaml:"-"`
	ItemsIds       []string                `json:"items" yaml:"items"`
	Effects        *[]ActiveEffectDocument `json:"effects" yaml:"effects"`
	Folder         string                  `json:"folder" yaml:"folder"`
	Sort           int                     `json:"sort" yaml:"sort"`
	Ownership      *Ownership              `json:"ownership" yaml:"ownership"`
	Flags          *Flags                  `json:"flags" yaml:"flags"`
	Stats          *DocumentStats          `json:"_stats" yaml:"_stats"`
}

func (d *ActorDocument) HydrateCollections(db *leveldb.DB) {}
