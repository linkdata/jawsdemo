package main

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/linkdata/jaws"
)

type uiClientPos struct{ *Globals }

func (ui uiClientPos) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	var sb strings.Builder
	ui.mu.RLock()
	for k := range ui.clientX {
		fmt.Fprintf(&sb, "%.0fx%.0f\n", ui.clientX[k], ui.clientY[k])
	}
	v = template.HTML(html.EscapeString(sb.String())) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (g *Globals) ClientPos() jaws.HtmlGetter {
	return uiClientPos{g}
}
