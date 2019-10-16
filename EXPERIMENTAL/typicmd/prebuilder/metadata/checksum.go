package metadata

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Checksum of file
func Checksum(target string) (updated bool, err error) {
	targetFilename := wd + "/" + target
	data, err := checksum(targetFilename)
	if err != nil {
		return
	}
	metaFilename := Path(filepath.Base(target) + ".md5")
	cachedData, err := ioutil.ReadFile(metaFilename)
	if os.IsNotExist(err) {
		updated = true
	} else {
		updated = (bytes.Compare(data, cachedData) != 0)
	}
	err = ioutil.WriteFile(metaFilename, data, 0777)
	return
}

func checksum(target string) (result []byte, err error) {
	file, err := os.Open(target)
	if err != nil {
		return
	}
	defer file.Close()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	result = hash.Sum(nil)
	result = []byte(hex.EncodeToString(result))
	return
}
