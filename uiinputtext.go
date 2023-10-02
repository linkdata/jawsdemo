package main

import (
	"fmt"

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
	ui.SelectPet.ReadLocked(func(nbl []*jaws.NamedBool) {
		for _, nb := range nbl {
			if nb.Name() == v {
				ui.SelectPet.Set(v, true)
				e.Jaws.Dirty(ui.SelectPet)
				break
			}
		}
	})
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
