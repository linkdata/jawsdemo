package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiCarsLink struct{ *Globals }

// 		{{$.A .G.CarsLinkID .G.CarsLinkText nil `href="/cars" class="btn btn-primary"`}}

func (ui *uiCarsLink) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.A(ui.CarsLinkID(), html.EscapeString(ui.CarsLinkText()), nil, attrs...)
}
