package main

import (
	"fmt"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputCheckbox struct{ *Globals }

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputCheckbox) eventFn(rq *jaws.Request, val bool) (err error) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if ui.InputCheckbox != val {
		ui.InputCheckbox = val
		rq.SetBoolValue(ui.InputCheckboxID(), val)
	}
	return
}

func (ui *uiInputCheckbox) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.Checkbox(ui.InputCheckboxID(), ui.InputCheckbox, ui.eventFn, attrs...) +
		template.HTML(fmt.Sprintf(`<label class="form-check-label" for="%s">Checkbox</label>`, ui.InputCheckboxID()))

}

func (uis *UiState) UiInputCheckbox() jaws.Ui { return &uiInputCheckbox{uis.G} }
