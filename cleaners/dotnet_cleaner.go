package cleaners

import (
	"time"

	"github.com/guionardo/dev-clean/internal"
)

type DotNetCleaner struct {
}

func (c *DotNetCleaner) Name() string {
	return "DotNet"
}

func (c *DotNetCleaner) CanCleanThisFolder(folderName string) bool {
	//TODO: Implementar canCleanThisFolder
	return false
}

func (c *DotNetCleaner) GetFiles(basePath string) (subFolders []string, fileNames []string, lastModTime time.Time, err error) {
	//TODO: Implementar getFiles
	return internal.GetFiles(basePath, "*.pyc")
}


func init() {
	AllCleaners().Register(&DotNetCleaner{})
}
