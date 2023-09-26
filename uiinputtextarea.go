package main

import (
	"github.com/linkdata/jaws"
)

type uiInputTextArea struct{ *Globals }

func (ui uiInputTextArea) JawsGetString(e *jaws.Element) (v string) {
	ui.mu.RLock()
	v = ui.inputTextArea
	ui.mu.RUnlock()
	return
}

func (ui uiInputTextArea) JawsSetString(e *jaws.Element, v string) (err error) {
	ui.mu.Lock()
	ui.inputTextArea = v
	ui.mu.Unlock()
	return
}

func (g *Globals) InputTextArea() jaws.StringSetter {
	return uiInputTextArea{g}
}
