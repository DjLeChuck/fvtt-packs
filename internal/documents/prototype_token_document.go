package documents

type PrototypeTokenDocument struct {
	baseTokenDocument `yaml:",inline"`
	Name              string `json:"name" yaml:"name"`
	RandomImg         bool   `json:"randomImg" yaml:"randomImg"`
}
