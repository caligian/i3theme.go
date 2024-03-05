package statusline

import (
	"fmt"
	"github.com/caligian/i3theme/config"
	"regexp"
	"strings"
)

//////////////////////////////////////////////////
type themeMap map[string][]string

type Statusline struct {
	Config   *config.Config
	Pos      int
	Colors   themeMap
	Contents []string
	Parsed   []string
	Spaces   int
}

//////////////////////////////////////////////////
var classes = [3]string{
	"border",
	"background",
	"separator",
}

var tabs = regexp.MustCompile("^[ ]*")
var whitespace = regexp.MustCompile("[ ]+")
var startRe = regexp.MustCompile("^[ ]*colors[ ]*[{]")
var endRe = regexp.MustCompile("^[ ]*[}][ ]*$")
var defaultTheme = themeMap{
	"background":         {"#000000"},
	"statusline":         {"#efefef"},
	"focused_workspace":  {"#ff0000", "#000000", "#ffffff"},
	"inactive_workspace": {"#282a2e", "#282a2e", "#ffffff"},
	"active_workspace":   {"#1d1f21", "#1d1f21", "#c5c8c6"},
	"urgent_workspace":   {"#2f343a", "#900000", "#ffffff"},
	"binding_mode":       {"#d2d5d3", "#000000", "#ffffff"},
}

//////////////////////////////////////////////////
func New(conf *config.Config) *Statusline {
	return &Statusline{
		Config:   conf,
		Contents: []string{},
		Colors:   themeMap{},
		Parsed:   []string{},
		Spaces:   8,
	}
}

func (st *Statusline) String() string {
	return strings.Join(st.Parsed, "\n")
}

func (st *Statusline) Read() *Statusline {
	conf := st.Config
	confContents := conf.Contents
	contents := st.Contents
	procLine := func(s string) string {
		line := whitespace.Split(s, -1)
		line = line[1:]
		name := line[0]
		rest := config.CheckHex(line[1:])

		if t := tabs.FindStringIndex(s); t != nil {
			st.Spaces = t[1] + 1
		}

		_, ok := defaultTheme[name]
		if !ok {
			panic(fmt.Sprintf(
				"invalid statusline color form: %s %s\n",
				name,
				strings.Join(rest, " "),
			))
		}

		st.Colors[name] = rest
		return s
	}

	for i, v := range confContents {
		if startRe.FindStringIndex(v) != nil {
			st.Pos = i + 1
			break
		}
	}

	if st.Pos == -1 {
		st.Colors = defaultTheme
		return st
	}

	for i := st.Pos; i < len(confContents); i++ {
		if ok := endRe.FindStringIndex(confContents[i]); ok != nil {
			break
		}
		contents = append(contents, procLine(confContents[i]))
	}

	st.Contents = contents
	return st
}

func (st *Statusline) Parse() *Statusline {
	parsed := st.Parsed
	for k, v := range defaultTheme {
		if _, ok := st.Colors[k]; !ok {
			st.Colors[k] = v
		}
	}

	for k, v := range st.Colors {
		parsed = append(parsed, fmt.Sprintf(
			"%s%-20s %s",
			strings.Repeat(" ", st.Spaces),
			k,
			strings.Join(v, " "),
		))
	}

	st.Parsed = parsed
	return st
}

func (st *Statusline) Sub() *config.Config {
	return st.Config.Sub(st.Pos, st.Contents, st.Parsed)
}

func Do(conf *config.Config) *config.Config {
  st := New(conf)
  return st.Read().Parse().Sub()
}
