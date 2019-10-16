package metadata

import (
	"fmt"
	"os"
)

var (
	metadataDir string
	wd          string
)

func init() {
	wd, _ = os.Getwd()
	metadataDir = wd + "/.typical-metadata"
	os.Mkdir(metadataDir, 0777)
}

// Path of metadata
func Path(filename string) string {
	return fmt.Sprintf("%s/%s", metadataDir, filename)
}

// CleanAll the metadata
func CleanAll() error {
	return os.RemoveAll(metadataDir)
}
