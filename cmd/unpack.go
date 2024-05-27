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
	"errors"
	"fmt"
	"github.com/djlechuck/fvtt-packs/internal/documents"
	"github.com/djlechuck/fvtt-packs/internal/fvttdb"
	"github.com/djlechuck/fvtt-packs/internal/serializer"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"os"
	"path/filepath"
	"strings"
)

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
			return fmt.Errorf("no directory \"%s\" found\n", d)
		}
		if err != nil {
			return fmt.Errorf("cannot access directory \"%s\": %s\n", d, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("\"%s\" is not a directory\n", d)
		}

		packs, err := os.ReadDir(pd)
		if err != nil {
			return fmt.Errorf("cannot read directory \"%s\": %s\n", pd, err)
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

			db, err := fvttdb.Open(fullPath)
			if err != nil {
				return fmt.Errorf("cannot open db: %s\n", err)
			}
			defer db.Close()

			err = db.IterateAll(func(iter iterator.Iterator) error {
				k := iter.Key()
				kStr := string(k)
				parts := strings.Split(kStr, "!")
				if len(parts) < 3 {
					return nil
				}
				collection := parts[1]
				if strings.Contains(collection, ".") {
					return nil // This is not a primary document, skip it.
				}

				fmt.Println("processing", kStr)
				doc, err := documents.Create(pName, collection, iter.Value())
				if err != nil {
					return fmt.Errorf("cannot get doc: %s\n", err)
				}

				if err := (*doc).HydrateCollections(db); err != nil {
					return fmt.Errorf("cannot hydrate doc collections: %s\n", err)
				}

				err = serializer.SerializeDocument(doc, "_pack_sources/"+pName, isYaml)
				if err != nil {
					return fmt.Errorf("cannot serialize doc: %s\n", err)
				}

				return nil
			})
			if err != nil {
				return fmt.Errorf("iterator error: %s\n", err)
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
