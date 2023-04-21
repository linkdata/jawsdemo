package main

import (
	"html/template"
	"strconv"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiInputRange struct {
	jid  string
	mu   deadlock.RWMutex // protects following
	data int
}

func newUiInputRange(jid string) *uiInputRange {
	return &uiInputRange{jid: jid}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputRange) eventFn(rq *jaws.Request, floatval float64) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	val := int(floatval)
	if ui.data != val {
		ui.data = val
		rq.SetFloatValue(ui.jid, floatval)
		rq.Jaws.SetInner(ui.jid+"-text", strconv.Itoa(ui.data))
	}
	return nil
}

func (ui *uiInputRange) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.Range(ui.jid, float64(ui.data), ui.eventFn, attrs...) +
		rq.Span(ui.jid+"-text", strconv.Itoa(ui.data), nil)
}
