package main

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/what"
)

type Globals struct {
	mu               deadlock.RWMutex
	InputText        *atomic.Value
	InputTextArea    string
	InputCheckbox    *atomic.Value
	InputRadioGroup1 *uiInputRadioGroup
	InputRadioGroup2 *uiInputRadioGroup
	InputDate        *atomic.Value
	InputRange       *atomic.Value
	InputButton      *uiInputButton
	SelectPet        *uiSelectPet
	Cars             []*Car
}

func NewGlobals() *Globals {
	g := &Globals{
		InputText:     &atomic.Value{},
		InputCheckbox: &atomic.Value{},
		InputRadioGroup1: newUiInputRadioGroup(
			jaws.NewNamedBoolArray("radiogroup1").Add("1", "Radio 1.1").Add("2", "Radio 1.2")),
		InputRadioGroup2: newUiInputRadioGroup(
			jaws.NewNamedBoolArray("radiogroup2").Add("1", "Radio 2.1").Add("2", "Radio 2.2")),
		InputDate:   &atomic.Value{},
		InputRange:  &atomic.Value{},
		InputButton: newUiInputButton(uiInputButtonID, "Meh"),
		SelectPet:   newUiSelectPet("selectpet"),
		Cars: []*Car{
			{
				VIN:       "JH4DB1671PS002584",
				Year:      1993,
				Make:      "Acura",
				Model:     "Integra",
				Condition: 56,
			},
			{
				VIN:       "KM8JT3AC2DU583865",
				Year:      2013,
				Make:      "Hyundai",
				Model:     "Tucson",
				Condition: 77,
			},
			{
				VIN:       "1D4GP24R75B188657",
				Year:      2005,
				Make:      "Dodge",
				Model:     "Grand Caravan",
				Condition: 89,
			},
		},
	}
	g.InputText.Store("")
	g.InputCheckbox.Store(false)
	g.InputDate.Store(time.Now())
	g.InputRange.Store(0.0)
	return g
}

func (g *Globals) RLock() {
	g.mu.RLock()
}

func (g *Globals) RUnlock() {
	g.mu.RUnlock()
}

func (g *Globals) SetInputButtonID() string {
	return "setinputbutton"
}

func (g *Globals) OnSetInputButton() jaws.EventFn {
	return func(rq *jaws.Request, evt what.What, id, val string) error {
		if evt == what.Trigger {
			g.mu.RLock()
			isWoo := strings.HasPrefix(g.InputButton.get(), "Woo")
			g.mu.RUnlock()
			if isWoo {
				rq.RemoveAttr(uiInputButtonID, "disabled")
			} else {
				rq.SetAttr(uiInputButtonID, "disabled", "")
			}
			g.mu.Lock()
			defer g.mu.Unlock()
			g.InputButton.set(val)
			rq.Jaws.SetInner(uiInputButtonID, g.InputButton.get())
		}
		return nil
	}
}

func (g *Globals) InputTextAreaID() string {
	return "inputtextarea"
}

func (g *Globals) OnInputTextArea() jaws.InputTextFn {
	return func(rq *jaws.Request, jid, val string) error {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.InputTextArea = val
		rq.SetTextValue(g.InputTextAreaID(), val)
		return nil
	}
}

func (g *Globals) InputRangeID() string {
	return "inputrange"
}
