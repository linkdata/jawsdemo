package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiCarsLink struct{ *Globals }

func (ui uiCarsLink) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	v = template.HTML(html.EscapeString(ui.carsLink)) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (g *Globals) CarsLink() jaws.HTMLGetter {
	return uiCarsLink{g}
}
