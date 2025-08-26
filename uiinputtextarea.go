package main

import (
	"github.com/linkdata/jaws"
)

func (g *Globals) InputTextArea() any {
	return jaws.Bind(&g.mu, &g.inputTextArea)
}
