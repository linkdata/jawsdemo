package main

import (
	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiInputText struct {
	mu   deadlock.RWMutex // protects following
	data string
}

func newUiInputText(jid, data string) jaws.UI {
	return jaws.NewUiText(jaws.ProcessTags(jid), &uiInputText{data: data})
}

func (ui *uiInputText) JawsGet(e *jaws.Element) (val interface{}) {
	ui.mu.RLock()
	val = ui.data
	ui.mu.RUnlock()
	return
}

func (ui *uiInputText) JawsSet(e *jaws.Element, val interface{}) (err error) {
	if s, ok := val.(string); ok {
		ui.mu.Lock()
		ui.data = s
		ui.mu.Unlock()
		for _, tag := range e.UI().JawsTags(e.Request()) {
			for _, elem := range e.Request().GetElements(tag) {
				e.Request().Jaws.SetValue(elem.Jid(), s)
			}
		}
	}
	return
}
