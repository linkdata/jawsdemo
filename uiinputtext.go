package main

import (
	"html"
	"html/template"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiInputText struct {
	jid  string
	mu   deadlock.RWMutex // protects following
	data string
}

func newUiInputText(jid, data string) *uiInputText {
	return &uiInputText{
		jid:  jid,
		data: data,
	}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputText) eventFn(rq *jaws.Request, val string) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	// it's usually a good idea to ensure that the value is actually changed before doing work
	if ui.data != val {
		ui.data = val
		// sends the changed text to all *other* Requests.
		rq.SetTextValue(ui.jid, val)
	}
	return nil
}

func (ui *uiInputText) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	ui.mu.RLock()
	data := ui.data
	ui.mu.RUnlock()
	return rq.Text(ui.jid, html.EscapeString(data), ui.eventFn, attrs...)
}
