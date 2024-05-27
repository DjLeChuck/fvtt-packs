package documents

import (
	"encoding/json"
	"fmt"
	"github.com/djlechuck/fvtt-packs/internal/fvttdb"
)

type ActorDocument struct {
	baseDocument   `yaml:",inline"`
	Img            string                  `json:"img" yaml:"img"`
	Type           string                  `json:"type" yaml:"type"`
	System         *System                 `json:"system" yaml:"system"`
	PrototypeToken *PrototypeTokenDocument `json:"prototypeToken" yaml:"prototypeToken"`
	Items          []*Document             `json:"-" yaml:"-"`
	ItemsIds       []string                `json:"-" yaml:"-"`
	Effects        *[]ActiveEffectDocument `json:"effects" yaml:"effects"`
	Folder         string                  `json:"folder" yaml:"folder"`
	Sort           int                     `json:"sort" yaml:"sort"`
	Ownership      *Ownership              `json:"ownership" yaml:"ownership"`
	Flags          *Flags                  `json:"flags" yaml:"flags"`
	Stats          *DocumentStats          `json:"_stats" yaml:"_stats"`
}

func (a *ActorDocument) UnmarshalJSON(data []byte) error {
	var aux struct {
		Items []string `json:"items"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	a.ItemsIds = aux.Items

	return nil
}

func (a *ActorDocument) MarshalJSON() ([]byte, error) {
	var aux struct {
		Items []string `json:"items"`
	}
	aux.Items = a.ItemsIds
	return json.Marshal(aux)
}

func (a *ActorDocument) HydrateCollections(fvttdb *fvttdb.FvttDb) error {
	fmt.Println(a.Name, a.ItemsIds)

	for _, id := range a.ItemsIds {
		key := "!actors.items!" + a.Id + "." + id
		v, err := fvttdb.Get(key)
		if err != nil {
			return fmt.Errorf("cannot get doc %s: %s\n", id, err)
		}

		doc, err := Create("", "items", v)
		if err != nil {
			return fmt.Errorf("cannot create doc %s: %s\n", id, err)
		}

		//fmt.Println(doc)
		a.Items = append(a.Items, doc)
	}

	return nil
}
