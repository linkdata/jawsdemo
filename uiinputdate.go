package main

import "github.com/linkdata/jaws/lib/bind"

func (g *Globals) InputDate() any {
	return bind.New(&g.mu, &g.inputDate)
}
