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
	"github.com/djlechuck/fvtt-packs/internal/documents"
	"github.com/djlechuck/fvtt-packs/internal/serializer"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path/filepath"
	"strings"
)

var documentTypeMapping = map[string]func() documents.Document{
	"actors":  func() documents.Document { return &documents.ActorDocument{} },
	"folders": func() documents.Document { return &documents.FolderDocument{} },
	"items":   func() documents.Document { return &documents.ItemDocument{} },
}

func getDoc(pack string, docType string, v []byte) (*documents.Document, error) {
	constructor, ok := documentTypeMapping[docType]
	if !ok {
		return nil, fmt.Errorf("structure not found for type %s", docType)
	}
	doc := constructor()
	if err := json.Unmarshal(v, &doc); err != nil {
		return nil, fmt.Errorf("cannot map document data: %s", err)
	}

	doc.SetPack(pack)
	doc.SetKey(docType)

	return &doc, nil
}

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack LevelDB into human-readable files",
	Long: `Unpack the LevelDB to get human-readable files. Multiple output formats are supported:
* JSON (default)
* YAML (with -y flag)

By default, packs are inside a packs directory. If this is not the case, you can override it with the -d flag: fvtt-packs unpack -d mypacks`,
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
			pName := pack.Name()
			if !pack.IsDir() {
				fmt.Println(pName, "is not a directory")
				continue
			}

			fmt.Println("unpacking", pName, "...")

			fullPath := filepath.Join(pd, pack.Name())
			db, err := leveldb.OpenFile(fullPath, &opt.Options{
				ErrorIfMissing: true,
				ReadOnly:       true,
			})
			if err != nil {
				fmt.Printf("cannot open pack \"%s\": %s\n", pName, err)
				continue
			}
			defer db.Close()

			iter := db.NewIterator(nil, nil)
			for iter.Next() {
				k := iter.Key()
				kStr := string(k)
				parts := strings.Split(kStr, "!")
				if len(parts) < 3 {
					continue
				}
				collection := parts[1]
				if strings.Contains(collection, ".") {
					continue // This is not a primary document, skip it.
				}

				fmt.Println("processing", kStr)
				doc, err := getDoc(pName, collection, iter.Value())
				if err != nil {
					fmt.Printf("cannot get doc: %s\n", err)
					continue
				}

				err = serializer.SerializeDocument(doc, "_pack_sources/"+pName, isYaml)
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
