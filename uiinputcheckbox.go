package main

import (
	"fmt"
	"html/template"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiInputCheckbox struct {
	jid  string
	mu   deadlock.RWMutex // protects following
	data bool
}

func newUiInputCheckbox(jid string) *uiInputCheckbox {
	return &uiInputCheckbox{jid: jid}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputCheckbox) eventFn(rq *jaws.Request, jid string, val bool) (err error) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if ui.data != val {
		ui.data = val
		rq.SetBoolValue(ui.jid, val)
	}
	return
}

func (ui *uiInputCheckbox) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	ui.mu.RLock()
	data := ui.data
	ui.mu.RUnlock()
	return rq.Checkbox(ui.jid, data, ui.eventFn, attrs...) +
		template.HTML(fmt.Sprintf(`<label class="form-check-label" for="%s">Checkbox</label>`, ui.jid))
}
