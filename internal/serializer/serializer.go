package serializer

import (
	"encoding/json"
	"github.com/djlechuck/fvtt-packs/internal/documents"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

func SerializeDocument(doc *documents.Document, destination string, isYaml bool) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
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

	err = os.WriteFile(path.Join(destination, (*doc).ExportName(isYaml)), append(serialized, '\n'), 0644)
	if err != nil {
		return err
	}

	return nil
}
