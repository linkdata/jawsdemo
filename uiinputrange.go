package main

import (
	"html/template"
	"strconv"

	"github.com/linkdata/jaws"
)

type uiInputRange struct{ *Globals }

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputRange) eventFn(rq *jaws.Request, floatval float64) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	val := int(floatval)
	if ui.InputRange != val {
		ui.InputRange = val
		rq.SetFloatValue(ui.InputRangeID(), floatval)
		// by using rq.Jaws.SetInner() we send the text to *all* Requests
		rq.Jaws.SetInner(ui.InputRangeTextID(), strconv.Itoa(ui.InputRange))
	}
	return nil
}

/*
	{{$.Range .G.InputRangeID (float64 .G.InputRange) .G.OnInputRange `class="form-range align-bottom w-75"`}}
	{{$.Span .G.InputRangeTextID (print .G.InputRange) nil ""}}
*/

func (ui *uiInputRange) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.Range(ui.InputRangeID(), float64(ui.InputRange), ui.eventFn, attrs...) +
		rq.Span(ui.InputRangeTextID(), strconv.Itoa(ui.InputRange), nil)
}
