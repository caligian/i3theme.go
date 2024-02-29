package main

import (
  "fmt"
  "github.com/caligian/i3theme/client"
  "github.com/caligian/i3theme/config"
)

type Config config.Config
type Client config.Client
type Statusline config.Statusline

func main() {
	config := config.New()
  fmt.Println("%#v\n", config)
}
