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
	for _, c := range ui.client {
		fmt.Fprintf(&sb, "%.0fx%.0f-%b ", c.X, c.Y, int(c.B))
	}
	ui.mu.RUnlock()
	v = template.HTML(html.EscapeString(sb.String())) //#nosec G203
	return
}

func (g *Globals) ClientPos() jaws.HtmlGetter {
	return uiClientPos{g}
}
