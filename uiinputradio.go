package main

import (
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputRadioGroup struct {
	nba *jaws.NamedBoolArray
}

func newUiInputRadioGroup(nba *jaws.NamedBoolArray) *uiInputRadioGroup {
	return &uiInputRadioGroup{
		nba: nba,
	}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputRadioGroup) eventFn(rq *jaws.Request, val string) error {
	ui.nba.SetOnly(val)
	rq.SetBoolValue(ui.nba.JidOf(val), true)
	return nil
}

func (ui *uiInputRadioGroup) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	// `<div class="form-check">%s<label class="form-check-label" for="%s">%s</label></div>`
	return rq.LabeledRadioGroup(ui.nba, ui.eventFn, attrs, []string{`class="form-check-label"`})
}
