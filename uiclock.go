package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiClock struct{ *Globals }

func (ui *uiClock) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.Div(ui.ClockID(), html.EscapeString(ui.ClockString()), nil, attrs...)
}
