package yizzy

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jimschubert/yizzy/models"
)

func TestMigrationDocParse(t *testing.T) {
	tests := []struct {
		name     string
		contents []byte
		wantErr  bool
		want     models.MigrationDoc
	}{
		{name: "full.json", contents: helperTestData(t, "operations", "full.json"), want: models.MigrationDoc{
			Operations: &[]models.MigrationEntry{
				{Op: "modify", Selector: strPtr(".ship-to.name"), Value: strPtr("env(FIRST_NAME)"), ValueType: strPtr("!!str")},
				{Op: "modify", Selector: strPtr(".total"), Value: strPtr("50")},
				{Op: "modify", Selector: strPtr(".modified-by"), Value: strPtr("tool")},
			},
			Env: &models.Env{
				"FIRST_NAME": ".bill-to.given",
				"LAST_NAME":  ".bill-to.family",
			},
		}},
		{name: "full.yaml", contents: helperTestData(t, "operations", "full.yaml"), want: models.MigrationDoc{
			Operations: &[]models.MigrationEntry{
				{Op: "modify", Selector: strPtr(".ship-to.name"), Value: strPtr("env(FIRST_NAME)"), ValueType: strPtr("!!str")},
				{Op: "modify", Selector: strPtr(".total"), Value: strPtr("50")},
				{Op: "modify", Selector: strPtr(".modified-by"), Value: strPtr("tool")},
			},
			Env: &models.Env{
				"FIRST_NAME": ".bill-to.given",
				"LAST_NAME":  ".bill-to.family",
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extension := strings.SplitAfter(tt.name, ".")
			fullPath, cleanup := tempFileCopy(t, tt.contents, extension[0])
			defer cleanup()

			doc := models.MigrationDoc{}
			err := doc.Load(fullPath)
			if tt.wantErr {
				assert.Error(t, err, "Load() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if !reflect.DeepEqual(doc, tt.want) {
					t.Errorf("doc.Load = %v, want %v", doc, tt.want)
				}
			}
		})
	}
}

