package main

import (
	"fmt"
	"strings"

	"github.com/linkdata/jaws"
)

type uiInputText struct{ *Globals }

func (ui uiInputText) JawsGetString(e *jaws.Element) (v string) {
	ui.mu.RLock()
	v = ui.inputText
	ui.mu.RUnlock()
	return
}

func (ui uiInputText) JawsSetString(e *jaws.Element, v string) (err error) {
	if v != "" && ui.SelectPet.Set(v, true) {
		e.Dirty(ui.SelectPet)
	}
	ui.mu.Lock()
	if strings.HasPrefix(v, "fail") {
		if v == "fail" {
			err = fmt.Errorf("whaddayamean, fail?")
		} else {
			// try using cut'n'paste or just holding down 'l'
			ui.inputText = "well, if you insist..."
		}
	} else {
		ui.inputText = v
	}
	ui.mu.Unlock()
	return
}

func (g *Globals) InputText() jaws.StringSetter {
	return uiInputText{g}
}
