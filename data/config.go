package data

import (
	"errors"
	"fmt"
	"os"
	"sylmark/lsp"
	"time"
)

type Config struct {
	RootMarkers *[]string `yaml:"root-markers" json:"rootMarkers"`
	RootPath    string
	DateLayout  string
}

func NewConfig() Config {
	rmakers := []string{".sylroot"}
	return Config{
		RootMarkers: &rmakers,
		DateLayout:  time.DateOnly,
	}
}

func (c *Config) GetDateString(date time.Time) string {
	return date.Format(c.DateLayout)
}
func (c *Config) CreatDirsIfNeeded() {
	c.CheckDirCreateIfNeeded("journal/")
}

func (c *Config) CheckDirCreateIfNeeded(dir string) (dirPath string, err error) {

	dirPath = fmt.Sprintf("%s/%s", c.RootPath, dir)
	stat, err := os.Stat(dirPath)
	if errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return dirPath, err
		}
	}

	return dirPath, nil
}

func (c *Config) GetFileURI(fileName string, dir string) (uri lsp.DocumentURI, error error) {
	urlPath := fmt.Sprintf("%s/%s%s", c.RootPath, dir, fileName)
	return UriFromPath(urlPath)

}
