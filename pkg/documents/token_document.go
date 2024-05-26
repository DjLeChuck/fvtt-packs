package documents

type TokenDocument struct {
	baseDocument `yaml:",inline"`
	ActorId      string      `json:"actorId" yaml:"actorId"`
	Delta        interface{} `json:"delta" yaml:"delta"`
	X            int         `json:"x" yaml:"x"`
	Y            int         `json:"y" yaml:"y"`
	Elevation    float64     `json:"elevation" yaml:"elevation"`
	Sort         int         `json:"sort" yaml:"sort"`
	Hidden       bool        `json:"hidden" yaml:"hidden"`
}
