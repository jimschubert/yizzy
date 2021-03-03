package yizzy

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

// helperTestData reads in test data by filename, combining all name parts and prefixing with 'testdata'
func helperTestData(t *testing.T, name ...string) []byte {
	t.Helper()
	pathParts := append([]string{"testdata"}, name...)
	path := filepath.Join(pathParts...) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

// strPtr turns string input into a pointer to that string
func strPtr(input string) *string {
	return &input
}

// hash returns a 32-bit FNV-1a hash of an input string. Helps avoid potential collisions on file naming.
func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// tempFileCopy copies test data contents to a hashed temp file, which helps avoid overwriting test files in place or
// working up previously manipulated test file contents.
func tempFileCopy(t *testing.T, data []byte, extension string) (fileLocation string, cleanup func()) {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "yizzy")
	if err != nil {
		t.Fatal(err)
	}
	r := rand.Int()
	testHash := hash(t.Name())
	testFile := fmt.Sprintf("file-%d-%d.%s", r, testHash, extension)
	filePath := filepath.Join(tempDir, testFile)
	if err := ioutil.WriteFile(filePath, data, 0600); err != nil {
		t.Fatal(err)
	}
	return filePath, func() { _ = os.RemoveAll(filePath) }
}
