package main

import (
	"fmt"
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiClientPos struct{ *Globals }

func (ui uiClientPos) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	s := fmt.Sprintf("%.0fx%.0f", ui.clientX, ui.clientY)
	v = template.HTML(html.EscapeString(s)) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (g *Globals) ClientPos() jaws.HtmlGetter {
	return uiClientPos{g}
}
