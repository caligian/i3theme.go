package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type themeMap map[string][]string

type Config struct {
	Path     string
	Contents []string
	Len      int
}

func (conf *Config) Read() *Config {
	dat, err := os.ReadFile(conf.Path)
	if err != nil {
		panic(err)
	}
	conf.Contents = strings.Split(string(dat), "\n")
	conf.Len = len(conf.Contents)

	return conf
}

func (conf *Config) Write() {
	s := conf.Contents
	err := os.WriteFile(conf.Path, []byte(strings.Join(s, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func CheckHex(s []string) []string {
	hex := []string{}
	hex_re := regexp.MustCompile("#[0-9a-fA-F]{6}")

	for i := 0; i < len(s); i++ {
		if len(s[i]) == 0 {
			continue
		}
		m := hex_re.FindString(s[i])
		if len(m) == 0 {
			panic(errors.New(fmt.Sprintf("invalid hex: %s", m)))
		}
		hex = append(hex, s[i])
	}

	return hex
}

func New() *Config {
	return &Config{
		Path: filepath.Join(os.Getenv("HOME"), ".config", "i3", "config"),
	}
}

func (conf *Config) String() string {
	return strings.Join(conf.Contents, "\n")
}

func (conf *Config) Sub(pos int, contents []string, parsed []string) *Config {
	confContents := conf.Contents
	contentsLen := len(contents)
	parsedLen := len(parsed)
	parsedHasMore := parsedLen > contentsLen

	if parsedHasMore {
		moreElems := parsed[contentsLen-1:]
		for i := 0; i < contentsLen-1; i++ {
			confContents[pos+i] = parsed[i]
		}
		confContents[pos+contentsLen-1] = strings.Join(moreElems, "\n")
	} else {
		for i, v := range parsed {
			confContents[pos+i] = v
		}
	}

	return conf
}

func Merge(defaultColors themeMap, colors themeMap) {
	for k, v := range colors {
		_, ok := defaultColors[k]
		if !ok {
			panic(fmt.Sprintf("invalid client theme form for %s %s\n", k, v))
		}
		colors[k] = CheckHex(v)
	}

  for k, v := range defaultColors {
    _, ok := colors[k]
    if !ok {
      colors[k] = v
    }
  }
}
