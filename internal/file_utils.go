package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/danwakefield/fnmatch"
)

func GetFiles(basePath string, masks ...string) (subFolders []string, fileNames []string, lastModTime time.Time, err error) {
	folders := make(map[string]interface{})

	basePath, err = filepath.Abs(basePath)
	if err != nil {
		return
	}
	lastModTime = time.Time{}

	err = filepath.Walk(basePath,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if lastModTime.Before(info.ModTime()) {
				lastModTime = info.ModTime()
			}
			if info.IsDir() {
				return nil
			}
			for _, mask := range masks {
				if fnmatch.Match(mask, info.Name(), fnmatch.FNM_FILE_NAME|fnmatch.FNM_IGNORECASE) {
					fileNames = append(fileNames, p)
					if path.Dir(p) != basePath {
						folders[path.Dir(p)] = true
						break
					}
				}
			}

			return nil
		})
	if err != nil {
		return
	}

	subFolders = make([]string, len(folders))
	i := 0
	for subFolder := range folders {
		subFolders[i] = subFolder
		i++
	}

	subFolders = sort.StringSlice(subFolders)
	fileNames = sort.StringSlice(fileNames)

	return

}

func GetFirstFile(basePath string, masks ...string) (fileName string, err error) {
	basePath, err = filepath.Abs(basePath)
	if err != nil {
		return
	}
	fileName = ""

	err = filepath.Walk(basePath,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			for _, mask := range masks {
				if fnmatch.Match(mask, info.Name(), fnmatch.FNM_FILE_NAME|fnmatch.FNM_IGNORECASE) {
					fileName = p
					return errors.New("Found")
				}
			}

			return nil
		})
	if err != nil && err.Error() == "Found" {
		err = nil
	}
	if len(fileName) == 0 {
		err = errors.New("Not found")
	}

	return

}

func ByteSizeToString(size int64) string {
	var suffix string
	v := float32(size)
	if size >= 1073741824 {
		v = v / 1073741824
		suffix = "GB"
	} else if size >= 1048576 {
		v = v / 1048576
		suffix = "MB"
	} else if size >= 1024 {
		v = v / 1024
		suffix = "KB"
	} else {
		return fmt.Sprintf("%d B", size)
	}
	return fmt.Sprintf("%.2f %s", v, suffix)

}
