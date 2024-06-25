package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiCarsLink struct{ *Globals }

func (ui uiCarsLink) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	v = template.HTML(html.EscapeString(ui.carsLink)) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (g *Globals) CarsLink() jaws.HtmlGetter {
	return uiCarsLink{g}
}
