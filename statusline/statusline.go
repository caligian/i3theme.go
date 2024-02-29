package statusline

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var StatuslineColorClasses = [3]string{
	"border",
	"background",
	"separator",
}



func (config *Config) ReadConfig(p string) *Config {
	lines, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}

	config.lines = strings.Split(string(lines), "\n")
	config.start = 0
	config.end = len(config.lines)

	return config
}

func (config *Config) ReadDefaultConfig() *Config {
	return config.ReadConfig(ConfigPath)
}

func (config *Config) ExtractStatusline() *Config {
	var status = Statusline{
		colors: map[string][]string{},
		lines:  []string{},
		start:  -1,
		end:    -1,
	}
	var lines = config.lines
	var end = len(lines)
  // needs fixing
	var status_re *regexp.Regexp = regexp.MustCompile("^status[.]([^ ]+) ([^$]+)")

	for i := 0; i < end; i++ {
		l := lines[i]
		matches := status_re.FindStringSubmatch(l)

		if matches == nil {
			continue
		} else if status.start == -1 {
			status.start = i
		}

		status_type := matches[1]
		colors := checkHex(strings.Split(matches[2], " "))
		status.colors[status_type] = colors
	}

	status.end = status.start + len(status.colors)
	for i := status.start; i < status.end; i++ {
		status.lines = append(status.lines, lines[i])
	}

	config.statusline = &status
	return config
}

func (config *Config) ExtractClient() *Config {
	var client = Client{
		colors: map[string][]string{},
		lines:  []string{},
		start:  -1,
		end:    -1,
	}
	var lines = config.lines
	var end = len(lines)
	var client_re *regexp.Regexp = regexp.MustCompile("^client[.]([^ ]+) ([^$]+)")

	for i := 0; i < end; i++ {
		l := lines[i]
		matches := client_re.FindStringSubmatch(l)

		if matches == nil {
			continue
		} else if client.start == -1 {
			client.start = i
		}

		client_type := matches[1]
		colors := checkHex(strings.Split(matches[2], " "))
		client.colors[client_type] = colors
	}

	client.end = client.start + len(client.colors)
	for i := client.start; i < client.end; i++ {
		client.lines = append(client.lines, lines[i])
	}

	config.client = &client
	return config
}
