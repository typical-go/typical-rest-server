package typiast

import (
	"io/ioutil"
)

// AllDirectories return all directory and sub-directory in path
func AllDirectories(path string, directories *[]string) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.IsDir() {
			dirPath := path + "/" + f.Name()
			AllDirectories(dirPath, directories)
			*directories = append(*directories, dirPath)
		}
	}
	return
}
