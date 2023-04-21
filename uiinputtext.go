package main

import (
	"html"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputText struct{ *Globals }

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputText) eventFn(rq *jaws.Request, val string) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	// it's usually a good idea to ensure that the value is actually changed before doing work
	if ui.InputText != val {
		ui.InputText = val
		// sends the changed text to all *other* Requests.
		rq.SetTextValue(ui.InputTextID(), val)
	}
	return nil
}

func (ui *uiInputText) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.Text(ui.InputTextID(), html.EscapeString(ui.InputText), ui.eventFn, attrs...)
}

func (uis *UiState) UiInputText() jaws.Ui { return &uiInputText{uis.G} }
