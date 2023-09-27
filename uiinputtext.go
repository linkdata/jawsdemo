package main

import (
	"fmt"

	"github.com/linkdata/jaws"
)

type uiInputText struct{ *Globals }

func (ui uiInputText) JawsGetString(rq *jaws.Request) (v string) {
	ui.mu.RLock()
	v = ui.inputText
	ui.mu.RUnlock()
	return
}

func (ui uiInputText) JawsSetString(rq *jaws.Request, v string) (err error) {
	ui.mu.Lock()
	if v == "fail" {
		err = fmt.Errorf("whaddayamean, fail?")
	} else {
		ui.inputText = v
	}
	ui.mu.Unlock()
	return
}

func (g *Globals) InputText() jaws.StringSetter {
	return uiInputText{g}
}
