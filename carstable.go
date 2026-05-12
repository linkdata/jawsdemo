package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/tag"
	"github.com/linkdata/jaws/lib/ui"
)

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
		tl = append(tl, ui.NewTemplate("car_row.html", c))
	}
	tl = append(tl, ui.NewTemplate("car_row.html", nil))
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
