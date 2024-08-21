package main

import (
	"github.com/linkdata/jaws"
)

type uiClientX struct{ *Globals }

func (ui uiClientX) JawsGetFloat(e *jaws.Element) (v float64) {
	ui.mu.RLock()
	v = ui.clientX[e.Session()]
	ui.mu.RUnlock()
	return
}

func (ui uiClientX) JawsSetFloat(e *jaws.Element, v float64) error {
	ui.mu.Lock()
	ui.clientX[e.Session()] = v
	ui.mu.Unlock()
	return nil
}

func (g *Globals) ClientX() jaws.FloatSetter {
	return uiClientX{g}
}
