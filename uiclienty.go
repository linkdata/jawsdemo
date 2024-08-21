package main

import (
	"github.com/linkdata/jaws"
)

type uiClientY struct{ *Globals }

func (ui uiClientY) JawsGetFloat(e *jaws.Element) (v float64) {
	ui.mu.RLock()
	v = ui.clientY[e.Session()]
	ui.mu.RUnlock()
	return
}

func (ui uiClientY) JawsSetFloat(e *jaws.Element, v float64) error {
	ui.mu.Lock()
	ui.clientY[e.Session()] = v
	ui.mu.Unlock()
	return nil
}

func (g *Globals) ClientY() jaws.FloatSetter {
	return uiClientY{g}
}
