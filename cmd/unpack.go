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
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

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
		subfolderPath := filepath.Join(p, d)

		info, err := os.Stat(subfolderPath)
		if os.IsNotExist(err) {
			return fmt.Errorf("no directory \"%s\" found", d)
		}
		if err != nil {
			return fmt.Errorf("cannot access directory \"%s\": %s", d, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("\"%s\" is not a directory", d)
		}

		fmt.Printf("coucou")

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
}
