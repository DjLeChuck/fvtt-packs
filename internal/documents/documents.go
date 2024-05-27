package documents

import (
	"encoding/json"
	"fmt"
	"github.com/djlechuck/fvtt-packs/internal/fvttdb"
	"regexp"
)

type Flags map[string]interface{}

type Ownership map[string]int

type System map[string]interface{}

type DocumentStats struct {
	CoreVersion      string `json:"coreVersion" yaml:"coreVersion"`
	SystemId         string `json:"systemId" yaml:"systemId"`
	SystemVersion    string `json:"systemVersion" yaml:"systemVersion"`
	CreatedTime      int    `json:"createdTime" yaml:"createdTime"`
	ModifiedTime     int    `json:"modifiedTime" yaml:"modifiedTime"`
	LastModifiedBy   string `json:"lastModifiedBy" yaml:"lastModifiedBy"`
	CompendiumSource string `json:"compendiumSource" yaml:"compendiumSource"`
	DuplicateSource  string `json:"duplicateSource" yaml:"duplicateSource"`
}

type TextureData struct {
	Src            string  `json:"src" yaml:"src"`
	AnchorX        float64 `json:"anchorX" yaml:"anchorX"`
	AnchorY        float64 `json:"anchorY" yaml:"anchorY"`
	OffsetX        float64 `json:"offsetX" yaml:"offsetX"`
	OffsetY        float64 `json:"offsetY" yaml:"offsetY"`
	Fit            string  `json:"fit" yaml:"fit"`
	ScaleX         float64 `json:"scaleX" yaml:"scaleX"`
	ScaleY         float64 `json:"scaleY" yaml:"scaleY"`
	Rotation       float64 `json:"rotation" yaml:"rotation"`
	Tint           string  `json:"tint" yaml:"tint"`
	AlphaThreshold float64 `json:"alphaThreshold" yaml:"alphaThreshold"`
}

type baseDocument struct {
	Pack string `json:"-" yaml:"-"`
	Key  string `json:"_key" yaml:"_key"`
	Id   string `json:"_id" yaml:"_id"`
	Name string `json:"name" yaml:"name"`
}

type Document interface {
	SetPack(pack string)
	SetKey(collection string)
	ExportName(isYaml bool) string
	HydrateCollections(fvttdb *fvttdb.FvttDb) error
}

var documentTypeMapping = map[string]func() Document{
	"actors":  func() Document { return &ActorDocument{} },
	"folders": func() Document { return &FolderDocument{} },
	"items":   func() Document { return &ItemDocument{} },
}

func (b *baseDocument) safeFilename() string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9А-я]`)

	return reg.ReplaceAllString(b.Name, "_")
}

func (b *baseDocument) SetPack(pack string) {
	b.Pack = pack
}

func (b *baseDocument) SetKey(collection string) {
	b.Key = "!" + collection + "!" + b.Id
}

func (b *baseDocument) ExportName(isYaml bool) string {
	extension := "json"
	if isYaml {
		extension = "yml"
	}

	if b.Name != "" {
		return b.safeFilename() + "_" + b.Id + "." + extension
	}

	return b.Key + "." + extension
}

func Create(pack string, docType string, v []byte) (*Document, error) {
	constructor, ok := documentTypeMapping[docType]
	if !ok {
		return nil, fmt.Errorf("structure not found for type %s\n", docType)
	}
	doc := constructor()
	if err := json.Unmarshal(v, &doc); err != nil {
		return nil, fmt.Errorf("cannot map document data: %s\n", err)
	}

	doc.SetPack(pack)
	doc.SetKey(docType)

	return &doc, nil
}
