package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputButton struct{ *Globals }

func (ui uiInputButton) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	v = template.HTML(html.EscapeString(ui.inputButton)) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (g *Globals) InputButton() jaws.HTMLGetter {
	return uiInputButton{g}
}
