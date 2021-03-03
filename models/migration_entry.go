package models

import "strings"

// MigrationEntry is the structure allowed by user input, as an array of operations in MigrationDoc
type MigrationEntry struct {
	Op        string  `json:"op" yaml:"op"`
	Selector  *string `json:"selector" yaml:"selector"`
	Value     *string `json:"value" yaml:"value"`
	Eval      *string `json:"eval" yaml:"eval"`
	ValueType *string `json:"value_type" yaml:"value_type"`
}

func (m *MigrationEntry) GetOperation() string {
	if m.Op == "" {
		return "modify"
	}

	return strings.ToLower(m.Op)
}
