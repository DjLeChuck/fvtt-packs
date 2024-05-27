package documents

type CombattantDocument struct {
	baseDocument `yaml:",inline"`
	Type         string         `json:"type" yaml:"type"`
	System       *System        `json:"system" yaml:"system"`
	ActorId      string         `json:"actorId" yaml:"actorId"`
	TokenId      string         `json:"tokenId" yaml:"tokenId"`
	SceneId      string         `json:"sceneId" yaml:"sceneId"`
	Img          string         `json:"img" yaml:"img"`
	Initiative   int            `json:"initiative" yaml:"initiative"`
	Hidden       bool           `json:"hidden" yaml:"hidden"`
	Defeated     bool           `json:"defeated" yaml:"defeated"`
	Flags        *Flags         `json:"flags" yaml:"flags"`
	Stats        *DocumentStats `json:"_stats" yaml:"_stats"`
}
