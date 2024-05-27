package documents

type CombatDocument struct {
	Id         string                `json:"_id" yaml:"_id"`
	Type       string                `json:"type" yaml:"type"`
	System     *System               `json:"system" yaml:"system"`
	Scene      *SceneDocument        `json:"scene" yaml:"scene"`
	Combatants *[]CombattantDocument `json:"combatants" yaml:"combatants"`
	Active     bool                  `json:"active" yaml:"active"`
	Round      int                   `json:"round" yaml:"round"`
	Turn       int                   `json:"turn" yaml:"turn"`
	Sort       int                   `json:"sort" yaml:"sort"`
	Flags      *Flags                `json:"flags" yaml:"flags"`
	Stats      *DocumentStats        `json:"_stats" yaml:"_stats"`
}
