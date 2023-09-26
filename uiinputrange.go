package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/linkdata/jaws"
)

type uiInputRange struct{ *Globals }

var _ jaws.UI = (*uiInputRange)(nil)          // statically ensure we implement this interface
var _ jaws.FloatGetter = (*uiInputRange)(nil) // statically ensure we implement this interface

func (ui uiInputRange) JawsRender(e *jaws.Element, w io.Writer, params []interface{}) {
	jaws.NewUiSpan(ui).JawsRender(e, w, params)
	e.Jaws.Dirty(ui)
}

func (ui uiInputRange) JawsUpdate(u jaws.Updater) {
	ui.mu.RLock()
	val := ui.inputRange
	ui.mu.RUnlock()
	switch {
	case val < 50:
		u.SetAttr("style", "color:red")
	case ui.inputRange < 90:
		u.RemoveAttr("style")
	default:
		u.SetAttr("style", "color:green")
	}
	jaws.NewUiSpan(ui).JawsUpdate(u)
}

func (ui uiInputRange) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
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
