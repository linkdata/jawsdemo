package main

import "github.com/linkdata/jaws/lib/bind"

func (g *Globals) InputTextArea() any {
	return bind.New(&g.mu, &g.inputTextArea)
}
