package main

import (
	"fmt"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputRange struct {
	jaws.Binder[float64]
	*Globals
}

func (ui *uiInputRange) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	switch {
	case ui.inputRange < 50:
		e.SetAttr("style", "color:red")
	case ui.inputRange < 90:
		e.RemoveAttr("style")
	default:
		e.SetAttr("style", "color:green")
	}
	v = template.HTML(fmt.Sprint(ui.inputRange)) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (g *Globals) InputRange() any {
	return &uiInputRange{
		Binder:  jaws.Bind(&g.mu, &g.inputRange),
		Globals: g,
	}
}
