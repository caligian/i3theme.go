package client

import (
	"fmt"
	"github.com/caligian/i3theme/config"
	"regexp"
	"strings"
)

type themeMap map[string][]string

type Client struct {
	Config   *config.Config
	Pos      int
	Contents []string
	Classes  [5]string
	Colors   themeMap
	Parsed   []string
}

var defaultTheme = map[string][]string{
	"focused":          {"#a54242", "#a54242", "#ffffff", "#2e9ef4", "#a54242"},
	"focused_inactive": {"#333333", "#5f676a", "#ffffff", "#484e50", "#5f676a"},
	"unfocused":        {"#333333", "#222222", "#888888", "#292d2e", "#222222"},
	"urgent":           {"#2f343a", "#900000", "#ffffff", "#900000", "#900000"},
	"placeholder":      {"#000000", "#0c0c0c", "#ffffff", "#000000", "#0c0c0c"},
	"background":       {"#ffffff"},
}

var classes = [5]string{
	"border",
	"background",
	"text",
	"indicator",
	"child_border",
}

func New(conf *config.Config) *Client {
	return &Client{
		Config:   conf,
		Classes:  classes,
		Pos:      -1,
		Contents: []string{},
		Colors:   themeMap{},
		Parsed:   []string{},
	}
}

func (cl *Client) Parse() *Client {
	theme := cl.Colors
	if len(theme) == 0 {
		theme = defaultTheme
	}
	s := []string{}

	for k, v := range defaultTheme {
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

	cl.Parsed = s
	return cl
}

func (cl *Client) Read() *Client {
	index := [2]int{-1, -1}
	checkStr := regexp.MustCompile("^client[.]")
	conf := cl.Config

	for i, v := range conf.Contents {
		if checkStr.FindStringIndex(v) != nil {
			index[0] = i
			break
		}
	}

	if index[0] == -1 {
		cl.Colors = defaultTheme
		cl.Pos = len(conf.Contents) + 1
		return cl
	}

	for i := index[0]; i < conf.Len; i++ {
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

	cl.Contents = conf.Contents[index[0] : index[1]+1]
	parsed := cl.Colors
	whitespace := regexp.MustCompile("[ ]+")
	sub := regexp.MustCompile("^client[.]")

	for _, v := range cl.Contents {
		if len(v) == 0 {
			continue
		}
		tokens := whitespace.Split(v, -1)
		parsed[sub.ReplaceAllString(tokens[0], "")] = tokens[1:]
	}

	cl.Pos = index[0]
	return cl
}

func (cl *Client) Sub() *config.Config {
  return cl.Config.Sub(cl.Pos, cl.Contents, cl.Parsed)
}

func Do(conf *config.Config) *config.Config {
  return New(conf).Read().Parse().Sub()
}
