package cleaners

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/guionardo/dev-clean/internal"
)

type (
	Cleaner interface {
		Name() string
		CanCleanThisFolder(folderName string) bool
		GetFiles(basePath string) (subFolders []string, fileNames []string, lastModTime time.Time, err error)
	}
	Cleaners struct {
		cleaners map[string]Cleaner
		lock     sync.Mutex
	}
)

var cleaners = &Cleaners{
	cleaners: make(map[string]Cleaner),
}

func AllCleaners() *Cleaners {
	return cleaners
}

func (c *Cleaners) Register(cleaner Cleaner) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cleaners[strings.ToLower(cleaner.Name())] = cleaner
}

func (c *Cleaners) GetNames() (names []string) {
	names = make([]string, len(c.cleaners))
	index := 0
	for name := range c.cleaners {
		names[index] = name
		index++
	}
	return
}

func (c *Cleaners) GetCleanerForFolder(folderName string) Cleaner {
	for _, cleaner := range c.cleaners {
		if cleaner.CanCleanThisFolder(folderName) {
			return cleaner
		}
	}
	return nil
}

func (c *Cleaners) GetCleanerFromName(name string, folderName string) (cleaner Cleaner, err error) {
	name = strings.ToLower(name)
	ok := false
	if name == "auto" {
		cleaner = c.GetCleanerForFolder(folderName)
		if cleaner == nil {
			err = fmt.Errorf("there's no cleaner for folder %s", folderName)
		}
	} else if cleaner, ok = c.cleaners[name]; !ok {
		err = fmt.Errorf("invalid cleaner %s", name)
	}
	if err == nil && cleaner != nil {
		if !cleaner.CanCleanThisFolder(folderName) {
			err = fmt.Errorf("cleaner %s cannot clean folder %s", name, folderName)
			cleaner = nil
		}
	}
	return
}

func (c *Cleaners) RunClean(cleaner Cleaner, folderName string, minimumAge time.Duration) error {
	fmt.Printf("Running cleaning [%s] @ %s", cleaner.Name(), folderName)
	_, fileNames, lastModTime, err := cleaner.GetFiles(folderName)
	if err != nil {
		return err
	}
	if len(fileNames) == 0 {
		fmt.Printf("No files to clean in [%s]\n", folderName)
		return nil
	}
	if minimumAge > 0 && lastModTime.After(time.Now().Add(-minimumAge)) {
		fmt.Printf("Last modification in [%s] is [%v] - skipped due minimum age [%v]", folderName, lastModTime, minimumAge)
		return nil
	}
	result, err := internal.RemoveFiles(folderName, fileNames, minimumAge)
	result.Write(os.Stdout)
	return err
}
