package documents

type tokenBarData struct {
	Attribute string `json:"attribute" yaml:"attribute"`
}

type sightData struct {
	Enabled     bool    `json:"enabled" yaml:"enabled"`
	Range       float64 `json:"range" yaml:"range"`
	Angle       float64 `json:"angle" yaml:"angle"`
	VisionMode  string  `json:"visionMode" yaml:"visionMode"`
	Color       string  `json:"color" yaml:"color"`
	Attenuation float64 `json:"attenuation" yaml:"attenuation"`
	Brightness  float64 `json:"brightness" yaml:"brightness"`
	Saturation  float64 `json:"saturation" yaml:"saturation"`
	Contrast    float64 `json:"contrast" yaml:"contrast"`
}

type detectionModeData struct {
	Id      string  `json:"id" yaml:"id"`
	Enabled bool    `json:"enabled" yaml:"enabled"`
	Range   float64 `json:"range" yaml:"range"`
}

type occludableData struct {
	Radius float64 `json:"radius" yaml:"radius"`
}

type ringColorData struct {
	Ring       string `json:"ring" yaml:"ring"`
	Background string `json:"background" yaml:"background"`
}

type ringSubjectData struct {
	Scale   float64 `json:"scale" yaml:"scale"`
	Texture string  `json:"texture" yaml:"texture"`
}

type ringData struct {
	Enabled bool             `json:"enabled" yaml:"enabled"`
	Colors  *ringColorData   `json:"colors" yaml:"colors"`
	Effects int              `json:"effects" yaml:"effects"`
	Subject *ringSubjectData `json:"subject" yaml:"subject"`
}

type baseTokenDocument struct {
	DisplayName      int                  `json:"displayName" yaml:"displayName"`
	ActorLink        bool                 `json:"actorLink" yaml:"actorLink"`
	AppendNumber     bool                 `json:"appendNumber" yaml:"appendNumber"`
	PrependAdjective bool                 `json:"prependAdjective" yaml:"prependAdjective"`
	Width            float64              `json:"width" yaml:"width"`
	Height           float64              `json:"height" yaml:"height"`
	Texture          *TextureData         `json:"texture" yaml:"texture"`
	HexagonalShape   int                  `json:"hexagonalShape" yaml:"hexagonalShape"`
	Locked           bool                 `json:"locked" yaml:"locked"`
	LockRotation     bool                 `json:"lockRotation" yaml:"lockRotation"`
	Rotation         float64              `json:"rotation" yaml:"rotation"`
	Alpha            float64              `json:"alpha" yaml:"alpha"`
	Disposition      int                  `json:"disposition" yaml:"disposition"`
	DisplayBars      int                  `json:"displayBars" yaml:"displayBars"`
	Bar1             *tokenBarData        `json:"bar1" yaml:"bar1"`
	Bar2             *tokenBarData        `json:"bar2" yaml:"bar2"`
	Light            interface{}          `json:"light" yaml:"light"`
	Sight            *sightData           `json:"sight" yaml:"sight"`
	DetectionModes   *[]detectionModeData `json:"detectionModes" yaml:"detectionModes"`
	Occludable       *occludableData      `json:"occludable" yaml:"occludable"`
	Ring             *ringData            `json:"ring" yaml:"ring"`
	Flags            *Flags               `json:"flags" yaml:"flags"`
}
