package main

import (
	"time"

	"github.com/linkdata/jaws"
)

type uiInputDate struct{ *Globals }

func (ui uiInputDate) JawsGetTime(e *jaws.Element) (v time.Time) {
	ui.mu.RLock()
	v = ui.inputDate
	ui.mu.RUnlock()
	return
}

func (ui uiInputDate) JawsSetTime(e *jaws.Element, v time.Time) (err error) {
	ui.mu.Lock()
	ui.inputDate = v
	ui.mu.Unlock()
	return
}

func (g *Globals) InputDate() jaws.TimeSetter {
	return uiInputDate{g}
}
