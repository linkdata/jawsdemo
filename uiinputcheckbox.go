package main

import (
	"github.com/linkdata/jaws"
)

type uiInputCheckbox struct{ *Globals }

func (ui uiInputCheckbox) JawsGetBool(e *jaws.Element) (v bool) {
	ui.mu.RLock()
	v = ui.inputCheckbox
	ui.mu.RUnlock()
	return
}

func (ui uiInputCheckbox) JawsSetBool(e *jaws.Element, v bool) (err error) {
	ui.mu.Lock()
	ui.inputCheckbox = v
	ui.mu.Unlock()
	return
}

func (g *Globals) InputCheckbox() jaws.BoolSetter {
	return uiInputCheckbox{g}
}
