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
		e.Jaws.SetAttr(btn.InputButton(), "disabled", "")
	} else {
		e.Jaws.RemoveAttr(btn.InputButton(), "disabled")
	}
	return
}

func (g *Globals) InputButton() any {
	// Click events will be handled in Globals.JawsClick()
	return uiInputButton{g}
}
