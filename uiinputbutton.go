package main

import (
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputButton struct{ *Globals }

func (ui uiInputButton) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	v = template.HTML(ui.inputButton)
	ui.mu.RUnlock()
	return
}

func (g *Globals) InputButton() jaws.HtmlGetter {
	return uiInputButton{g}
}
