package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".config", "i3", "config")
var StatuslineColorClasses = [3]string{
	"border",
	"background",
	"separator",
}
var ClientColorClasses = [5]string{
	"border",
	"background",
	"text",
	"indicator",
	"child_border",
}


type Statusline struct {
	Start, End int
	Contents   []string
	Colors     map[string][]string
	// background, statusline, separator string
	// focused_workspace                 [3]string
	// active_workspace                  [3]string
	// inactive_workspace                [3]string
	// urgent_workspace                  [3]string
	// binding_mode                      [3]string
}


type Client struct {
	Start, End int
	Contents   []string
	Colors     map[string][]string
}

type Config struct {
	Contents   []string
	Client     *Client
	Statusline *Statusline
  Len int
}

func (statusline *Statusline) NewClient(contents []string) *Statusline {
	return &Statusline{
		Start:    -1,
		End:      -1,
		Contents: contents,
		Colors:   map[string][]string{},
	}
}

func ReadConfig() []string {
	config, err := os.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(config), "\n")
}

func WriteConfig(s []string) {
	err := os.WriteFile(ConfigPath, []byte(strings.Join(s, "\n")), 0644)
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
  return &Config{Contents: ReadConfig()}
}
