package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/linkdata/jaws"
)

type uiInputDate struct{ *Globals }

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputDate) eventFn(rq *jaws.Request, val time.Time) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	// it's usually a good idea to ensure that the value is actually changed before doing work
	if ui.InputDate != val {
		ui.InputDate = val
		// sends the changed text to all *other* Requests.
		rq.SetDateValue(ui.InputDateID(), val)
	}
	return nil
}

func (ui *uiInputDate) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return template.HTML(fmt.Sprintf(`<label for="%s" class="form-label">Date</label>`, ui.InputDateID())) +
		rq.Date(ui.InputDateID(), ui.InputDate, ui.eventFn, attrs...)
}
