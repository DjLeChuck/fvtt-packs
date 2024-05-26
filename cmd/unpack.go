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
	CoreVersion      string
	SystemId         string
	SystemVersion    string
	CreatedTime      int
	ModifiedTime     int
	LastModifiedBy   string
	CompendiumSource string
	DuplicateSource  string
}

type TextureData struct {
	Src            string
	AnchorX        float64
	AnchorY        float64
	OffsetX        float64
	OffsetY        float64
	Fit            string
	ScaleX         float64
	ScaleY         float64
	Rotation       float64
	Tint           string
	AlphaThreshold float64
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
	Key  string `json:"-"`
	Id   string `json:"_id"`
	Name string
}

type FolderDocument struct {
	BaseDocument
	Type        string
	Description string
	Folder      *FolderDocument
	Sorting     string
	Sort        int
	Color       string
	Flags       *Flags
	Stats       *DocumentStats `json:"_stats"`
}

type SceneDocument struct {
	BaseDocument
	Active              bool
	Navigation          bool
	NavOrder            int
	NavName             string
	Background          *TextureData
	Foreground          string
	ForegroundElevation int
	Thumb               string
	Width               int
	Height              int
	Padding             int
	Initial             struct {
		X     int
		Y     int
		Scale float64
	}
	BackgroundColor string
	Grid            struct {
		Type      int
		Size      int
		Style     string
		Thickness int
		Color     string
		Alpha     float64
		Distance  float64
		Units     string
	}
	TokenVision          bool
	FogExploration       bool
	FogReset             int
	GlobalLight          bool
	GlobalLightThreshold float64
	Darkness             float64
	FogOverlay           string
	FogExploredColor     string
	FogUnexploredColor   string
	Drawings             []interface{}
	Tokens               []interface{}
	Lights               []interface{}
	Notes                []interface{}
	Sounds               []interface{}
	Templates            []interface{}
	Tiles                []interface{}
	Walls                []interface{}
	Playlist             []interface{}
	PlaylistSound        []interface{}
	Journal              []interface{}
	JournalEntryPage     []interface{}
	Weather              string
	Folder               *FolderDocument
	Sort                 int
	Ownership            *Ownership
	Flags                *Flags
	Stats                *DocumentStats `json:"_stats"`
}

type CombattantDocument struct {
	BaseDocument
	Type       string
	System     *System
	ActorId    string
	TokenId    string
	SceneId    string
	Img        string
	Initiative int
	Hidden     bool
	Defeated   bool
	Flags      *Flags
	Stats      *DocumentStats `json:"_stats"`
}

type CombatDocument struct {
	Id         string `json:"_id"`
	Type       string
	System     *System
	Scene      *SceneDocument
	Combatants *[]CombattantDocument
	Active     bool
	Round      int
	Turn       int
	Sort       int
	Flags      *Flags
	Stats      *DocumentStats `json:"_stats"`
}

type ActiveEffectDocument struct {
	BaseDocument
	Img     string
	Type    string
	System  *System
	Changes []struct {
		Key      string
		Value    string
		Mode     int
		Priority float64
	}
	Disabled bool
	Duration []struct {
		StartTime  int
		Seconds    int
		Combat     *CombatDocument
		Rounds     int
		Turns      int
		StartRound int
		StartTurn  int
	}
	Description string
	Origin      string
	Tint        string
	Transfer    bool
	Statuses    []string
	Flags       *Flags
	Stats       *DocumentStats `json:"_stats"`
}

type ItemDocument struct {
	BaseDocument
	Type      string
	Img       string
	System    *System
	Effects   *[]ActiveEffectDocument
	Folder    *FolderDocument
	Sort      int
	Ownership *Ownership
	Flags     *Flags
	Stats     *DocumentStats `json:"_stats"`
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
				err = serializeDocument(doc, doc.ExportName(false), false)
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
}
