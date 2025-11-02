package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sylmark/lsp"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RootMarkers              []string
	IncludeMdExtensionMdLink bool `toml:"include_md_extension_md_link"`
	MdLinkWebMode            bool `toml:"md_link_web_mode"`
	RootPath                 string
	DateLayout               string
}

func NewConfig() Config {
	rmakers := []string{".sylroot.toml"}
	return Config{
		RootMarkers:              rmakers,
		IncludeMdExtensionMdLink: true,
		DateLayout:               time.DateOnly,
		MdLinkWebMode:            false,
	}
}

func (c *Config) LoadConfig() {
	filePath := filepath.Join(c.RootPath, ".sylroot.toml")
	toml.DecodeFile(filePath, c)
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

// could be dirInVault = "journal"
func (c *Config) GetFileURI(fileName string, dirInVault string) (uri lsp.DocumentURI, error error) {
	urlPath := filepath.Join(c.RootPath, dirInVault, fileName)
	return UriFromPath(urlPath)
}
