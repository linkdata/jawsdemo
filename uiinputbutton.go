package main

import (
	"html/template"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

const uiInputButtonID = "inputbutton"

type uiInputButton struct {
	jid  string
	mu   deadlock.RWMutex // protects following
	data string
}

func newUiInputButton(jid, data string) *uiInputButton {
	return &uiInputButton{
		jid:  jid,
		data: data,
	}
}

func (ui *uiInputButton) get() (data string) {
	ui.mu.RLock()
	data = ui.data
	ui.mu.RUnlock()
	return
}

func (ui *uiInputButton) set(data string) {
	ui.mu.Lock()
	ui.data = data
	ui.mu.Unlock()
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputButton) eventFn(rq *jaws.Request, jid string) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if ui.data != "Bar" {
		ui.data = "Bar"
	} else {
		rq.Alert("info", "Foo?")
		ui.data = "<strong>Foo</strong>"
	}
	rq.Jaws.SetInner(uiInputButtonID, ui.data)
	return nil
}

func (ui *uiInputButton) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	ui.mu.RLock()
	data := ui.data
	ui.mu.RUnlock()
	return rq.Button(ui.jid, data, ui.eventFn, attrs...)
}
