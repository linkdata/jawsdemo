package main

import (
	"github.com/linkdata/jaws"
)

func (g *Globals) InputCheckbox() any {
	return jaws.Bind(&g.mu, &g.inputCheckbox)
}
