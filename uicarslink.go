package main

import (
	"github.com/linkdata/jaws"
)

func (g *Globals) CarsLink() any {
	return jaws.Bind(&g.mu, &g.carsLink)
}
