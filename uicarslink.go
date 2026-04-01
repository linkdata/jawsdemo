package main

import "github.com/linkdata/jaws/lib/bind"

func (g *Globals) CarsLink() any {
	return bind.New(&g.mu, &g.carsLink)
}
