package main

import (
	"fmt"
	"html/template"

	"github.com/linkdata/jaws"
)

type uiInputRange struct{ *Globals }

var _ jaws.FloatGetter = (*uiInputRange)(nil) // statically ensure we implement this interface
var _ jaws.HtmlGetter = (*uiInputRange)(nil)  // statically ensure we implement this interface

func (ui uiInputRange) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	switch {
	case ui.inputRange < 50:
		e.SetAttr("style", "color:red")
	case ui.inputRange < 90:
		e.RemoveAttr("style")
	default:
		e.SetAttr("style", "color:green")
	}
	v = template.HTML(fmt.Sprint(ui.inputRange))
	ui.mu.RUnlock()
	return
}

func (ui uiInputRange) JawsGetFloat(e *jaws.Element) (v float64) {
	ui.mu.RLock()
	v = float64(ui.inputRange)
	ui.mu.RUnlock()
	return
}

func (ui uiInputRange) JawsSetFloat(e *jaws.Element, v float64) (err error) {
	ui.mu.Lock()
	ui.inputRange = int(v)
	ui.mu.Unlock()
	return
}

func (g *Globals) InputRange() interface{} {
	return uiInputRange{g}
}
