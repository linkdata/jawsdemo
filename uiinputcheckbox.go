package main

import (
	"github.com/linkdata/jaws"
)

type uiInputCheckbox struct{ *Globals }

func (ui uiInputCheckbox) JawsGetBool(e *jaws.Element) (v bool) {
	return ui.inputCheckbox.Load()
}

func (ui uiInputCheckbox) JawsSetBool(e *jaws.Element, v bool) (err error) {
	ui.inputCheckbox.Store(v)
	return
}

func (g *Globals) InputCheckbox() jaws.BoolSetter {
	return uiInputCheckbox{g}
}
