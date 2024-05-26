package documents

type SceneDocument struct {
	baseDocument        `yaml:",inline"`
	Active              bool         `json:"active" yaml:"active"`
	Navigation          bool         `json:"navigation" yaml:"navigation"`
	NavOrder            int          `json:"navOrder" yaml:"navOrder"`
	NavName             string       `json:"navName" yaml:"navName"`
	Background          *TextureData `json:"background" yaml:"background"`
	Foreground          string       `json:"foreground" yaml:"foreground"`
	ForegroundElevation int          `json:"foregroundElevation" yaml:"foregroundElevation"`
	Thumb               string       `json:"thumb" yaml:"thumb"`
	Width               int          `json:"width" yaml:"width"`
	Height              int          `json:"height" yaml:"height"`
	Padding             int          `json:"padding" yaml:"padding"`
	Initial             struct {
		X     int     `json:"x" yaml:"x"`
		Y     int     `json:"y" yaml:"y"`
		Scale float64 `json:"scale" yaml:"scale"`
	} `json:"initial" yaml:"initial"`
	BackgroundColor string `json:"backgroundColor" yaml:"backgroundColor"`
	Grid            struct {
		Type      int     `json:"type" yaml:"type"`
		Size      int     `json:"size" yaml:"size"`
		Style     string  `json:"style" yaml:"style"`
		Thickness int     `json:"thickness" yaml:"thickness"`
		Color     string  `json:"color" yaml:"color"`
		Alpha     float64 `json:"alpha" yaml:"alpha"`
		Distance  float64 `json:"distance" yaml:"distance"`
		Units     string  `json:"units" yaml:"units"`
	} `json:"grid" yaml:"grid"`
	TokenVision          bool           `json:"tokenVision" yaml:"tokenVision"`
	FogExploration       bool           `json:"fogExploration" yaml:"fogExploration"`
	FogReset             int            `json:"fogReset" yaml:"fogReset"`
	GlobalLight          bool           `json:"globalLight" yaml:"globalLight"`
	GlobalLightThreshold float64        `json:"globalLightThreshold" yaml:"globalLightThreshold"`
	Darkness             float64        `json:"darkness" yaml:"darkness"`
	FogOverlay           string         `json:"fogOverlay" yaml:"fogOverlay"`
	FogExploredColor     string         `json:"fogExploredColor" yaml:"fogExploredColor"`
	FogUnexploredColor   string         `json:"fogUnexploredColor" yaml:"fogUnexploredColor"`
	Drawings             []interface{}  `json:"drawings" yaml:"drawings"`
	Tokens               []interface{}  `json:"tokens" yaml:"tokens"`
	Lights               []interface{}  `json:"lights" yaml:"lights"`
	Notes                []interface{}  `json:"notes" yaml:"notes"`
	Sounds               []interface{}  `json:"sounds" yaml:"sounds"`
	Templates            []interface{}  `json:"templates" yaml:"templates"`
	Tiles                []interface{}  `json:"tiles" yaml:"tiles"`
	Walls                []interface{}  `json:"walls" yaml:"walls"`
	Playlist             []interface{}  `json:"playlist" yaml:"playlist"`
	PlaylistSound        []interface{}  `json:"playlistSound" yaml:"playlistSound"`
	Journal              []interface{}  `json:"journal" yaml:"journal"`
	JournalEntryPage     []interface{}  `json:"journalEntryPage" yaml:"journalEntryPage"`
	Weather              string         `json:"weather" yaml:"weather"`
	Folder               string         `json:"folder" yaml:"folder"`
	Sort                 int            `json:"sort" yaml:"sort"`
	Ownership            *Ownership     `json:"ownership" yaml:"ownership"`
	Flags                *Flags         `json:"flags" yaml:"flags"`
	Stats                *DocumentStats `json:"_stats" yaml:"_stats"`
}
