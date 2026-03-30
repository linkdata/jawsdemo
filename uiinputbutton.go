package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputButton struct{ *Globals }

func (btn uiInputButton) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	btn.mu.RLock()
	v = template.HTML(html.EscapeString(btn.inputButton)) //#nosec G203
	btn.mu.RUnlock()
	if e.Session().Get("mystical") != nil {
		e.SetAttr("disabled", "")
	} else {
		e.RemoveAttr("disabled")
	}
	return
}

func (g *Globals) InputButton() any {
	// Click events will be handled in Globals.JawsClick()
	return uiInputButton{g}
}
