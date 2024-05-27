package documents

type ActiveEffectDocument struct {
	baseDocument `yaml:",inline"`
	Img          string  `json:"img" yaml:"img"`
	Type         string  `json:"type" yaml:"type"`
	System       *System `json:"system" yaml:"system"`
	Changes      []struct {
		Key      string  `json:"key" yaml:"key"`
		Value    string  `json:"value" yaml:"value"`
		Mode     int     `json:"mode" yaml:"mode"`
		Priority float64 `json:"priority" yaml:"priority"`
	} `json:"changes" yaml:"changes"`
	Disabled bool `json:"disabled" yaml:"disabled"`
	Duration []struct {
		StartTime  int             `json:"startTime" yaml:"startTime"`
		Seconds    int             `json:"seconds" yaml:"seconds"`
		Combat     *CombatDocument `json:"combat" yaml:"combat"`
		Rounds     int             `json:"rounds" yaml:"rounds"`
		Turns      int             `json:"turns" yaml:"turns"`
		StartRound int             `json:"startRound" yaml:"startRound"`
		StartTurn  int             `json:"startTurn" yaml:"startTurn"`
	} `json:"duration" yaml:"duration"`
	Description string         `json:"description" yaml:"description"`
	Origin      string         `json:"origin" yaml:"origin"`
	Tint        string         `json:"tint" yaml:"tint"`
	Transfer    bool           `json:"transfer" yaml:"transfer"`
	Statuses    []string       `json:"statuses" yaml:"statuses"`
	Flags       *Flags         `json:"flags" yaml:"flags"`
	Stats       *DocumentStats `json:"_stats" yaml:"_stats"`
}
