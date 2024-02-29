package main

import (
	"fmt"
	"github.com/caligian/i3theme/config"
	"regexp"
	"strings"
)

type clientTheme map[string][]string

var DefaultTheme = clientTheme{
	"focused":          {"#a54242", "#a54242", "#ffffff", "#2e9ef4", "#a54242"},
	"focused_inactive": {"#333333", "#5f676a", "#ffffff", "#484e50", "#5f676a"},
	"unfocused":        {"#333333", "#222222", "#888888", "#292d2e", "#222222"},
	"urgent":           {"#2f343a", "#900000", "#ffffff", "#900000", "#900000"},
	"placeholder":      {"#000000", "#0c0c0c", "#ffffff", "#000000", "#0c0c0c"},
	"background":       {"#ffffff"},
}

func Parse(conf *config.Config) {
	theme := conf.Client.Colors
	if len(theme) == 0 {
		theme = DefaultTheme
	}

	s := []string{}

	for k, v := range DefaultTheme {
		_, ok := theme[k]
		if !ok {
			theme[k] = v
		}
	}

	for k, v := range theme {
		k = "client." + k
		colors := strings.Join(config.CheckHex(v), " ")
		s = append(s, fmt.Sprintf("%-25s %s", k, colors))
	}

	conf.Client.Contents = s
}

func New(conf *config.Config) {
	cl := &config.Client{
		Start:    -1,
		End:      -1,
		Contents: []string{},
		Colors:   map[string][]string{},
	}
	conf.Client = cl
}

func GetContents(conf *config.Config) {
	index := [2]int{-1, -1}
	checkStr := regexp.MustCompile("^client[.]")

	for i, v := range conf.Contents {
		if checkStr.FindStringIndex(v) != nil {
			index[0] = i
			break
		}
	}

	if index[0] == -1 {
		conf.Client.Colors = DefaultTheme
		conf.Client.Start = len(conf.Contents) + 1
		conf.Client.End = conf.Client.Start + len(conf.Client.Colors)
		return
	}

	for i := index[0]; i < len(conf.Contents); i++ {
		v := conf.Contents[i]
		if len(v) == 0 {
			continue
		} else if checkStr.FindStringIndex(v) != nil {
			index[1] = i
		} else {
			break
		}
	}

	if index[1] == -1 {
		index[1] = index[0] + 1
	}

	conf.Client.Contents = conf.Contents[index[0]:index[1]]
	parsed := conf.Client.Colors
	whitespace := regexp.MustCompile("[ ]+")
	sub := regexp.MustCompile("^client[.]")

	for _, v := range conf.Client.Contents {
		if len(v) == 0 {
			continue
		}
		tokens := whitespace.Split(v, -1)
		parsed[sub.ReplaceAllString(tokens[0], "")] = tokens[1:]
	}

	conf.Client.Start = index[0]
	conf.Client.End = index[1] + 1
}

func Sub(conf *config.Config) {
	contents := conf.Contents
	contentsLen := len(contents)
	cl := conf.Client
	clFoundLen := cl.End - cl.Start
	clLen := len(cl.Colors)
	clContents := cl.Contents

	if cl.Start > contentsLen {
		contents = append(contents, cl.Contents...)
	} else if clFoundLen < clLen {
		j := 0
		l := cl.Start + clFoundLen

		for i := cl.Start; i < cl.Start+clFoundLen; i++ {
			contents[i] = clContents[j]
			j = j + 1
		}

		if j < clLen {
			here := contents[l]
			contents[l] = strings.Join(append(clContents[j:], "", here), "\n")
		}
	} else {
		j := 0
		l := cl.Start + clLen

		for i := cl.Start; i < l; i++ {
			contents[i] = clContents[j]
			j = j + 1
		}

		for i := l; i < cl.End+clFoundLen; i++ {
			contents[i] = ""
		}
	}

	fmt.Printf("%s\n", strings.Join(contents, "\n"))
}

func main() {
	conf := config.New()
	// theme := &DefaultTheme
	New(conf)
	GetContents(conf)
	Parse(conf)
	Sub(conf)
}
