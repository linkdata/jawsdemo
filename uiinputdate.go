package main

import (
	"github.com/linkdata/jaws"
)

func (g *Globals) InputDate() any {
	return jaws.Bind(&g.mu, &g.inputDate)
}
