package yizzy

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jimschubert/yizzy/models"
)

// func TestYamlMigrator_Run(t *testing.T) {
// 	type fields struct {
// 		InputFile           string
// 		MigrationsDirectory string
// 		InPlace             bool
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			y := &YamlMigrator{
// 				InputFile:           tt.fields.InputFile,
// 				MigrationsDirectory: tt.fields.MigrationsDirectory,
// 				InPlace:             tt.fields.InPlace,
// 			}
// 			if err := y.Run(); (err != nil) != tt.wantErr {
// 				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestYamlMigrator_processFile(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "use_envs.yaml", wantErr: false},
		{name: "new_property.yaml", wantErr: false},
		{name: "modify_property.yaml", wantErr: false},
		{name: "add_comment.yaml", wantErr: false},
		{name: "alternative_selectors.yaml", wantErr: false},
	}

	for _, tt := range tests {
		input := tt.name
		expectOutput := "expect_" + input

		inputFile := helperTestData(t, "invoice.yaml")
		operations := helperTestData(t, "operations", input)
		expected := helperTestData(t, "expectations", expectOutput)

		extension := strings.SplitAfter(tt.name, ".")
		inputCopy, cleanInput := tempFileCopy(t, inputFile, extension[0])
		migrationFile, cleanup := tempFileCopy(t, operations, extension[0])

		teardown := func() {
			cleanup()
			cleanInput()
		}

		t.Run(tt.name, func(t *testing.T) {
			defer teardown()
			y := &YamlMigrator{
				InputFile: inputCopy,
				InPlace:   true,
			}

			doc := &models.MigrationDoc{}
			err := doc.Load(migrationFile)
			if tt.wantErr {
				assert.Error(t, err, "processFile() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				resetYqLogger()
				if err := y.processFile(doc); (err != nil) != tt.wantErr {
					t.Errorf("processFile() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					// file compare
					actual, err := ioutil.ReadFile(inputCopy)
					assert.NoError(t, err)
					assert.NotNil(t, actual)
					if err != nil {
						t.Errorf("processFile failed to load actual file")
					} else {
						assert.Equal(t, string(expected), string(actual))
					}
				}
			}
		})
	}
}

// func TestYamlMigrator_validateFile(t *testing.T) {
// 	type fields struct {
// 		InputFile           string
// 		MigrationsDirectory string
// 		InPlace             bool
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			y := &YamlMigrator{
// 				InputFile:           tt.fields.InputFile,
// 				MigrationsDirectory: tt.fields.MigrationsDirectory,
// 				InPlace:             tt.fields.InPlace,
// 			}
// 			if err := y.validateFile(); (err != nil) != tt.wantErr {
// 				t.Errorf("validateFile() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
//
// func TestYamlMigrator_validateMigrationsDirectory(t *testing.T) {
// 	type fields struct {
// 		InputFile           string
// 		MigrationsDirectory string
// 		InPlace             bool
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			y := &YamlMigrator{
// 				InputFile:           tt.fields.InputFile,
// 				MigrationsDirectory: tt.fields.MigrationsDirectory,
// 				InPlace:             tt.fields.InPlace,
// 			}
// 			if err := y.validateMigrationsDirectory(); (err != nil) != tt.wantErr {
// 				t.Errorf("validateMigrationsDirectory() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
