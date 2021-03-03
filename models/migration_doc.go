package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Env is a map of keys to selectors. For example:
// NAME: '.a.b.name'
type Env map[string]string

// MigrationDoc describes the operations to be applied, and environment mappings bound to these.
// An Env is a map of keys to bind as environment/context, and values which are document selectors to query the value.
type MigrationDoc struct {
	Operations *[]MigrationEntry `json:"operations" yaml:"operations"`
	Env        *Env              `json:"env" yaml:"env"`
}

// Load a MigrationDoc via path. Supports .json and .yaml/.yml files. Any file without a .json extension is parsed as YAML.
func (c *MigrationDoc) Load(path string) error {
	b, err := ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if strings.HasSuffix(path, ".json") {
		return json.Unmarshal(b, c)
	}

	return yaml.Unmarshal(b, c)
}
