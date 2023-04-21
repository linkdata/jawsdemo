package main

import (
	"strings"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type Globals struct {
	mu            deadlock.RWMutex
	InputText     *uiInputText
	InputTextArea string
	InputCheckbox *uiInputCheckbox
	InputRadio1   *uiInputRadio
	InputRadio2   *uiInputRadio
	InputDate     *uiInputDate
	InputRange    *uiInputRange
	InputButton   *uiInputButton
	SelectPet     *uiSelectPet
	Cars          []*Car
}

func NewGlobals() *Globals {
	return &Globals{
		InputText:     newUiInputText("inputtext", ""),
		InputCheckbox: newUiInputCheckbox("checkbox"),
		InputRadio1:   newUiInputRadio("inputradio1/a", "Radio 1", false),
		InputRadio2:   newUiInputRadio("inputradio2/a", "Radio 2", false),
		InputDate:     newUiInputDate("inputdate"),
		InputRange:    newUiInputRange("inputrange"),
		InputButton:   newUiInputButton(uiInputButtonID, "Meh"),
		SelectPet:     newUiSelectPet("selectpet"),
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
	return func(rq *jaws.Request, id, evt, val string) error {
		if evt == "trigger" {
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
	return func(rq *jaws.Request, val string) error {
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

func (g *Globals) InputRangeTextID() string {
	return "inputrangetext"
}

func (g *Globals) InputCheckboxID() string {
	return "inputcheckbox"
}
