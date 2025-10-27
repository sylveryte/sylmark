package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sylmark/lsp"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RootMarkers              []string
	IncludeMdExtensionMdLink bool `toml:"include_md_extension_md_link"`
	RootPath                 string
	DateLayout               string
}

func NewConfig() Config {
	rmakers := []string{"sylroot.toml"}
	return Config{
		RootMarkers:              rmakers,
		IncludeMdExtensionMdLink: true,
		DateLayout:               time.DateOnly,
	}
}

func (c *Config) LoadConfig() {
	filePath := filepath.Join(c.RootPath, "sylroot.toml")
	toml.DecodeFile(filePath, c)
}

// removes .md file if config demands
func (c *Config) GetMdFormattedTargetUrl(path string) string {
	if c.IncludeMdExtensionMdLink {
		return path
	}
	return RemoveMdExtOnly(path)
}

// adds .md where needed
func (c *Config) GetMdRealUrlAndSubTarget(fullUrl string) (url lsp.DocumentURI, subTarget SubTarget, found bool) {
	found = strings.ContainsRune(fullUrl, '#')
	if found {
		splits := strings.SplitN(fullUrl, "#", 2)
		url = lsp.DocumentURI(splits[0])
		subTarget = SubTarget("#" + splits[1])
	} else {
		url = lsp.DocumentURI(fullUrl)
	}
	url = lsp.DocumentURI(GetMdRealTargetUrl(string(url)))
	return url, subTarget, found
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
