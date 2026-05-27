package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/tag"
	"github.com/linkdata/jaws/lib/ui"
)

const carRowTemplate = "car_row.html"

type CarsTable struct {
	globals *Globals
}

var _ tag.TagGetter = (*CarsTable)(nil)

func (ct *CarsTable) owner() *Globals {
	if ct != nil && ct.globals != nil {
		return ct.globals
	}
	return globals
}

func (ct *CarsTable) JawsGetTag(tag.Context) any {
	return &ct.owner().Cars
}

func (ct *CarsTable) JawsContains(e *jaws.Element) (tl []jaws.UI) {
	g := ct.owner()
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, c := range g.Cars {
		tl = append(tl, CarRow{Car: c})
	}
	tl = append(tl, CarRow{})
	return tl
}

func (ct *CarsTable) JawsClick(e *jaws.Element, data jaws.Click) (err error) {
	switch data.Name {
	case "add":
		ct.owner().AddRandomCar()
		e.Dirty(ct)
		return nil
	}
	return jaws.ErrEventUnhandled
}

func (ct *CarsTable) Mystical() jaws.ClickHandler {
	return ui.New("Mystical").Clicked(func(obj ui.Object, elem *jaws.Element, click jaws.Click) (err error) {
		elem.Session().Set("mystical", nil)
		elem.Dirty(ct.owner().InputButton())
		return
	})
}

type CarRow struct {
	Car *Car
}

var _ jaws.UI = CarRow{}
var _ jaws.ClickHandler = CarRow{}

func (row CarRow) auth(elem *jaws.Element) jaws.Auth {
	if f := elem.Request.Jaws.MakeAuth; f != nil {
		return f(elem.Request)
	}
	return jaws.DefaultAuth{}
}

func (row CarRow) execute(elem *jaws.Element, w io.Writer) error {
	tmpl := elem.Request.Jaws.LookupTemplate(carRowTemplate)
	if tmpl == nil {
		return fmt.Errorf("missing template %q: %w", carRowTemplate, ui.ErrMissingTemplate)
	}
	return tmpl.Execute(w, ui.With{
		Element:       elem,
		RequestWriter: ui.RequestWriter{Request: elem.Request, Writer: w},
		Dot:           row.Car,
		Auth:          row.auth(elem),
	})
}

func (row CarRow) JawsRender(elem *jaws.Element, w io.Writer, params []any) (err error) {
	elem.Tag(row.Car)
	attrs := elem.ApplyParams(params)

	b := elem.Jid().AppendStartTagAttr(nil, "tr")
	for _, attr := range attrs {
		if attr != "" {
			b = append(b, ' ')
			b = append(b, attr...)
		}
	}
	b = append(b, '>')
	if _, err = w.Write(b); err != nil {
		return err
	}
	if err = row.execute(elem, w); err != nil {
		return err
	}
	_, err = io.WriteString(w, "</tr>")
	return err
}

func (row CarRow) JawsUpdate(elem *jaws.Element) {
	var sb strings.Builder
	if err := row.execute(elem, &sb); err == nil {
		elem.SetInner(template.HTML(sb.String())) // #nosec G203
	} else {
		elem.Request.MustLog(err)
	}
}

func (row CarRow) JawsClick(elem *jaws.Element, click jaws.Click) error {
	if row.Car == nil {
		return jaws.ErrEventUnhandled
	}
	return row.Car.JawsClick(elem, click)
}
