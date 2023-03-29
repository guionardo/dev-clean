package cleaners

import (
	"time"

	"github.com/guionardo/dev-clean/internal"
)

type PythonCleaner struct {
}

func (c *PythonCleaner) Name() string {
	return "Python"
}

func (c *PythonCleaner) CanCleanThisFolder(folderName string) bool {
	_, err := internal.GetFirstFile(folderName, "*.py")
	return err == nil
}

func (c *PythonCleaner) GetFiles(basePath string) (subFolders []string, fileNames []string, lastModTime time.Time, err error) {
	return internal.GetFiles(basePath, "*.pyc")
}



func init() {
	AllCleaners().Register(&PythonCleaner{})
}
