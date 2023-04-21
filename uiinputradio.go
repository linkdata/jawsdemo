package main

import (
	"fmt"
	"html/template"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiInputRadio struct {
	jid  string
	name string
	mu   deadlock.RWMutex // protects following
	data bool
}

func newUiInputRadio(jid, name string, data bool) *uiInputRadio {
	return &uiInputRadio{
		jid:  jid,
		name: name,
		data: data,
	}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputRadio) eventFn(rq *jaws.Request, val bool) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	// it's usually a good idea to ensure that the value is actually changed before doing work
	if ui.data != val {
		ui.data = val
		// sends the changed value to all *other* Requests.
		rq.SetBoolValue(ui.jid, val)
	}
	return nil
}

func (ui *uiInputRadio) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	ui.mu.RLock()
	data := ui.data
	ui.mu.RUnlock()
	return template.HTML(fmt.Sprintf(`<div class="form-check">%s<label class="form-check-label" for="%s">%s</label></div>`,
		rq.Radio(ui.jid, data, ui.eventFn, attrs...),
		ui.jid, ui.name,
	))
}
