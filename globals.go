package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type Globals struct {
	mu            deadlock.RWMutex
	InputText     *uiInputText
	InputTextArea string
	InputCheckbox bool
	InputRadio1   *uiInputRadio
	InputRadio2   *uiInputRadio
	InputDate     time.Time
	InputRange    int
	InputButton   string
	SelectPet     string
	Cars          []*Car
}

func NewGlobals() *Globals {
	return &Globals{
		InputText:   newUiInputText("inputtext", ""),
		InputRadio1: newUiInputRadio("inputradio1/a", "Radio 1", false),
		InputRadio2: newUiInputRadio("inputradio2/a", "Radio 2", false),
		InputButton: "Meh",
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
			isWoo := strings.HasPrefix(g.InputButton, "Woo")
			g.mu.RUnlock()
			if isWoo {
				rq.RemoveAttr(g.InputButtonID(), "disabled")
			} else {
				rq.SetAttr(g.InputButtonID(), "disabled", "")
			}
			g.mu.Lock()
			defer g.mu.Unlock()
			g.InputButton = val
			rq.Jaws.SetInner(g.InputButtonID(), g.InputButton)
		}
		return nil
	}
}

func (g *Globals) InputButtonID() string {
	return "inputbutton"
}

func (g *Globals) OnInputButton() jaws.ClickFn {
	return func(rq *jaws.Request) error {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.InputButton != "Bar" {
			g.InputButton = "Bar"
		} else {
			rq.Alert("info", "Foo?")
			g.InputButton = "<strong>Foo</strong>"
		}
		rq.Jaws.SetInner(g.InputButtonID(), g.InputButton)
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

func (g *Globals) OnInputCheckbox() jaws.InputBoolFn {
	return func(rq *jaws.Request, val bool) (err error) {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.InputCheckbox = val
		rq.SetBoolValue(g.InputCheckboxID(), val)
		return
	}
}

func (g *Globals) InputDateID() string {
	return "inputdate"
}

func (g *Globals) OnInputDate() jaws.InputDateFn {
	return func(rq *jaws.Request, val time.Time) (err error) {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.InputDate = val
		rq.SetDateValue(g.InputDateID(), val)
		return
	}
}

func (g *Globals) SelectPetID() string {
	return "selectpet"
}

func (g *Globals) OnSelectPet() jaws.InputTextFn {
	return func(rq *jaws.Request, val string) (err error) {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.SelectPet = val
		rq.SetTextValue(g.SelectPetID(), val)
		return
	}
}

func (g *Globals) SelectPetOptions() (sol *jaws.NamedBoolArray) {
	sol = jaws.NewNamedBoolArray()
	sol.Add("", "--Please choose an option--")
	sol.Add("dog", "Dog")
	sol.Add("cat", "Cat")
	sol.Add("hamster", "Hamster")
	sol.Add("parrot", "Parrot")
	sol.Add("spider", "Spider")
	sol.Check(g.SelectPet)
	return
}

func (g *Globals) CarsLinkID() string {
	return "cars"
}

func (g *Globals) CarsLinkText() string {
	switch rand.Intn(5) {
	case 0:
		return "Check out these cars!"
	case 1:
		return "Did you know VIN numbers are encoded?"
	case 2:
		return "DO NOT CLICK HERE!"
	case 3:
		return "Cars"
	}
	return "This is a boring link to car info."
}
