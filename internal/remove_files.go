package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type (
	FileInfo struct {
		Path      string
		Status    FileStatus
		Size      int64
		Deleted   bool
		TimeStamp time.Time
	}
	FileStatus byte

	RemovedFilesResult struct {
		BasePath     string
		RemovedFiles []string
		RemovedSize  int64
		TotalSize    int64
	}
)

const (
	NotFound FileStatus = iota
	Found
)

func GetFilesInfos(filesNames []string) (filesInfos []FileInfo) {
	filesInfos = make([]FileInfo, len(filesNames))

	for index, fileName := range filesNames {
		if stat, err := os.Stat(fileName); err == nil && !stat.IsDir() {
			filesInfos[index] = FileInfo{
				Path:      fileName,
				Status:    Found,
				Size:      stat.Size(),
				TimeStamp: stat.ModTime(),
			}
		} else {
			filesInfos[index] = FileInfo{
				Path:   fileName,
				Status: NotFound,
			}
		}
	}

	return
}

func RemoveFiles(basePath string, filesNames []string, minimumAge time.Duration) (result RemovedFilesResult, err error) {
	result = RemovedFilesResult{
		BasePath:     basePath,
		RemovedFiles: make([]string, 0, len(filesNames)),
		TotalSize:    0,
	}
	// folders := make(map[string]bool)
	filesInfos := GetFilesInfos(filesNames)
	for _, fileInfo := range filesInfos {
		if fileInfo.Status != Found || fileInfo.TimeStamp.After(time.Now().Add(-minimumAge)) {
			continue
		}
		result.TotalSize += fileInfo.Size

		if err := os.Remove(fileInfo.Path); err == nil {
			result.RemovedFiles = append(result.RemovedFiles, fileInfo.Path)
			result.RemovedSize += fileInfo.Size
			// folders[path.Dir(fileInfo.Path)] = true
		}
	}
	return
}

func (r *RemovedFilesResult) Write(w io.Writer) {
	const line = "----------------------------------"
	fmt.Fprintf(w, "Files removed from %s\n", r.BasePath)
	for _, file := range r.RemovedFiles {
		fmt.Fprintf(w, "\t%s\n", strings.TrimPrefix(file, r.BasePath))
	}
	fmt.Fprintf(w, `%s
Removed Size  %s
Removed Files %d
%s`, line, ByteSizeToString(r.RemovedSize), len(r.RemovedFiles), line)
}
