package main

import (
	"html/template"
	"io"
	"strconv"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/what"
)

type uiInputRange struct {
	mu      deadlock.Mutex
	data    float64
	uiRange *jaws.UiRange
}

func newUiInputRange(jid string) (ui *uiInputRange) {
	ui = &uiInputRange{}
	ui.uiRange = jaws.NewUiRange([]interface{}{jid}, ui)
	return
}

func (ui *uiInputRange) JawsGet(e *jaws.Element) (val interface{}) {
	ui.mu.Lock()
	val = ui.data
	ui.mu.Unlock()
	return

}
func (ui *uiInputRange) JawsSet(e *jaws.Element, val interface{}) (changed bool) {
	ui.mu.Lock()
	changed = val != ui.data
	ui.data = val.(float64)
	ui.mu.Unlock()
	return
}

func (ui *uiInputRange) JawsInner(e *jaws.Element) (s template.HTML) {
	ui.mu.Lock()
	s = template.HTML(ui.formatSpan())
	ui.mu.Unlock()
	return
}

func (ui *uiInputRange) formatSpan() string {
	return strconv.FormatFloat(ui.data, 'f', 1, 64)
}

func (ui *uiInputRange) JawsCreate(rq *jaws.Request, data []interface{}) (elem *jaws.Element, err error) {
	return rq.NewElement(nil, ui.uiRange, data), nil
}

func (ui *uiInputRange) JawsRender(e *jaws.Element, w io.Writer) (err error) {
	return ui.uiRange.JawsRender(e, w)
}

func (ui *uiInputRange) JawsUpdate(e *jaws.Element) (err error) {
	if err = ui.uiRange.JawsUpdate(e); err == nil {
		e.Request().Jaws.Update([]interface{}{ui})
	}
	return
}

func (ui *uiInputRange) JawsEvent(e *jaws.Element, wht what.What, val string) (err error) {
	return ui.uiRange.JawsEvent(e, wht, val)
}
