package prebuilder

import "io/ioutil"

func scanProject(root string) (dirs, files []string, err error) {
	dirs = append(dirs, root)
	err = scanDir(root, &dirs, &files)
	return
}

func scanDir(root string, directories, files *[]string) (err error) {
	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}
	for _, f := range fileInfos {
		if f.IsDir() {
			dirPath := root + "/" + f.Name()
			scanDir(dirPath, directories, files)
			*directories = append(*directories, dirPath)
		} else {
			*files = append(*files, root+"/"+f.Name())
		}
	}
	return
}
