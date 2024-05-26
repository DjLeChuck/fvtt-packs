/*
Copyright © 2024 DjLeChuck <djlechuck@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

type Document interface {
	SetKey(collection string)
	ExportName(isYaml bool) string
}

func (b *BaseDocument) SetKey(collection string) {
	b.Key = "!" + collection + "!" + b.Id
}

func getSafeFilename(filename string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9А-я]`)

	return reg.ReplaceAllString(filename, "_")
}

func (b *BaseDocument) ExportName(isYaml bool) string {
	extension := "json"
	if isYaml {
		extension = "yml"
	}

	if b.Name != "" {
		return getSafeFilename(b.Name) + "_" + b.Id + "." + extension
	}

	return b.Key + "." + extension
}

type BaseDocument struct {
	Key  string `json:"_key" yaml:"_key"`
	Id   string `json:"_id" yaml:"_id"`
	Name string `json:"name" yaml:"name"`
}

type FolderDocument struct {
	BaseDocument
	Type        string          `json:"type" yaml:"type"`
	Description string          `json:"description" yaml:"description"`
	Folder      *FolderDocument `json:"folder" yaml:"folder"`
	Sorting     string          `json:"sorting" yaml:"sorting"`
	Sort        int             `json:"sort" yaml:"sort"`
	Color       string          `json:"color" yaml:"color"`
	Flags       *Flags          `json:"flags" yaml:"flags"`
	Stats       *DocumentStats  `json:"_stats" yaml:"_stats"`
}

type SceneDocument struct {
	BaseDocument
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
	TokenVision          bool            `json:"tokenVision" yaml:"tokenVision"`
	FogExploration       bool            `json:"fogExploration" yaml:"fogExploration"`
	FogReset             int             `json:"fogReset" yaml:"fogReset"`
	GlobalLight          bool            `json:"globalLight" yaml:"globalLight"`
	GlobalLightThreshold float64         `json:"globalLightThreshold" yaml:"globalLightThreshold"`
	Darkness             float64         `json:"darkness" yaml:"darkness"`
	FogOverlay           string          `json:"fogOverlay" yaml:"fogOverlay"`
	FogExploredColor     string          `json:"fogExploredColor" yaml:"fogExploredColor"`
	FogUnexploredColor   string          `json:"fogUnexploredColor" yaml:"fogUnexploredColor"`
	Drawings             []interface{}   `json:"drawings" yaml:"drawings"`
	Tokens               []interface{}   `json:"tokens" yaml:"tokens"`
	Lights               []interface{}   `json:"lights" yaml:"lights"`
	Notes                []interface{}   `json:"notes" yaml:"notes"`
	Sounds               []interface{}   `json:"sounds" yaml:"sounds"`
	Templates            []interface{}   `json:"templates" yaml:"templates"`
	Tiles                []interface{}   `json:"tiles" yaml:"tiles"`
	Walls                []interface{}   `json:"walls" yaml:"walls"`
	Playlist             []interface{}   `json:"playlist" yaml:"playlist"`
	PlaylistSound        []interface{}   `json:"playlistSound" yaml:"playlistSound"`
	Journal              []interface{}   `json:"journal" yaml:"journal"`
	JournalEntryPage     []interface{}   `json:"journalEntryPage" yaml:"journalEntryPage"`
	Weather              string          `json:"weather" yaml:"weather"`
	Folder               *FolderDocument `json:"folder" yaml:"folder"`
	Sort                 int             `json:"sort" yaml:"sort"`
	Ownership            *Ownership      `json:"ownership" yaml:"ownership"`
	Flags                *Flags          `json:"flags" yaml:"flags"`
	Stats                *DocumentStats  `json:"_stats" yaml:"_stats"`
}

type CombattantDocument struct {
	BaseDocument
	Type       string         `json:"type" yaml:"type"`
	System     *System        `json:"system" yaml:"system"`
	ActorId    string         `json:"actorId" yaml:"actorId"`
	TokenId    string         `json:"tokenId" yaml:"tokenId"`
	SceneId    string         `json:"sceneId" yaml:"sceneId"`
	Img        string         `json:"img" yaml:"img"`
	Initiative int            `json:"initiative" yaml:"initiative"`
	Hidden     bool           `json:"hidden" yaml:"hidden"`
	Defeated   bool           `json:"defeated" yaml:"defeated"`
	Flags      *Flags         `json:"flags" yaml:"flags"`
	Stats      *DocumentStats `json:"_stats" yaml:"_stats"`
}

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

type ActiveEffectDocument struct {
	BaseDocument
	Img     string  `json:"img" yaml:"img"`
	Type    string  `json:"type" yaml:"type"`
	System  *System `json:"system" yaml:"system"`
	Changes []struct {
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

type ItemDocument struct {
	BaseDocument
	Type      string                  `json:"type" yaml:"type"`
	Img       string                  `json:"img" yaml:"img"`
	System    *System                 `json:"system" yaml:"system"`
	Effects   *[]ActiveEffectDocument `json:"effects" yaml:"effects"`
	Folder    *FolderDocument         `json:"folder" yaml:"folder"`
	Sort      int                     `json:"sort" yaml:"sort"`
	Ownership *Ownership              `json:"ownership" yaml:"ownership"`
	Flags     *Flags                  `json:"flags" yaml:"flags"`
	Stats     *DocumentStats          `json:"_stats" yaml:"_stats"`
}

var documentTypeMapping = map[string]func() Document{
	"items":   func() Document { return &ItemDocument{} },
	"folders": func() Document { return &FolderDocument{} },
}

func getDoc(docType string, v []byte) (Document, error) {
	constructor, ok := documentTypeMapping[docType]
	if !ok {
		return nil, fmt.Errorf("structure not found for type %s", docType)
	}
	doc := constructor()
	if err := json.Unmarshal(v, &doc); err != nil {
		return nil, fmt.Errorf("cannot map document data: %s", err)
	}

	return doc, nil
}

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack LevelDB into human readable files",
	Long: `Unpack the LevelDB to get human readable files. Multiple output formats are supported:
* JSON (default)
* YAML (with -y flag)

By default, packs are inside a packs directory. If this is not the case, you can override it with the -d flag: foundrypacks unpack -d mypacks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, _ := cmd.Flags().GetString("path")
		if p == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return errors.New("cannot get the current working directory")
			}
			p = cwd
		}

		d, _ := cmd.Flags().GetString("directory")
		pd := filepath.Join(p, d)

		info, err := os.Stat(pd)
		if os.IsNotExist(err) {
			return fmt.Errorf("no directory \"%s\" found", d)
		}
		if err != nil {
			return fmt.Errorf("cannot access directory \"%s\": %s", d, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("\"%s\" is not a directory", d)
		}

		packs, err := os.ReadDir(pd)
		if err != nil {
			return fmt.Errorf("cannot read directory \"%s\": %s", pd, err)
		}

		isYaml, _ := cmd.Flags().GetBool("yaml")

		for _, pack := range packs {
			if !pack.IsDir() {
				fmt.Println(pack.Name(), "is not a directory")
				continue
			}

			fmt.Println("unpacking", pack.Name())

			fullPath := filepath.Join(pd, pack.Name())
			db, err := leveldb.OpenFile(fullPath, &opt.Options{
				ErrorIfMissing: true,
				ReadOnly:       true,
			})
			if err != nil {
				fmt.Printf("cannot open pack \"%s\": %s\n", pack.Name(), err)
				continue
			}
			defer db.Close()

			iter := db.NewIterator(nil, nil)
			for iter.Next() {
				k := iter.Key()
				parts := strings.Split(string(k), "!")
				if len(parts) < 3 {
					continue
				}
				collection := parts[1]
				if strings.Contains(collection, ".") {
					continue // This is not a primary document, skip it.
				}

				doc, err := getDoc(collection, iter.Value())
				if err != nil {
					fmt.Printf("cannot get doc: %s\n", err)
					continue
				}

				doc.SetKey(collection)
				err = serializeDocument(doc, doc.ExportName(isYaml), isYaml)
				if err != nil {
					fmt.Printf("cannot serialize doc: %s\n", err)
				}
			}
			iter.Release()
			if err := iter.Error(); err != nil {
				return fmt.Errorf("iterator error: %s", err)
			}
		}

		return nil
	},
}

func serializeDocument(doc interface{}, filename string, isYaml bool) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	var serialized []byte
	var err error

	if isYaml {
		serialized, err = yaml.Marshal(doc)
		if err != nil {
			return err
		}
	} else {
		serialized, err = json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(filename, append(serialized, '\n'), 0644)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(unpackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unpackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	unpackCmd.Flags().StringP("path", "p", "", "Path of the directory containing LevelDB packs")
	unpackCmd.Flags().StringP("directory", "d", "packs", "Directory containing LevelDB packs")
	unpackCmd.Flags().BoolP("yaml", "y", false, "Unpack as YAML files instead of JSON")
}
