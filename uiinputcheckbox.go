package main

import "github.com/linkdata/jaws/lib/bind"

func (g *Globals) InputCheckbox() any {
	return bind.New(&g.mu, &g.inputCheckbox)
}
