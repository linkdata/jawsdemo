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
	inputText        string
	inputTextArea    string
	InputCheckbox    *atomic.Value
	InputRadioGroup1 *jaws.NamedBoolArray
	InputRadioGroup2 *jaws.NamedBoolArray
	InputDate        *atomic.Value
	InputRange       *atomic.Value
	InputButton      *atomic.Value
	SelectPet        *jaws.NamedBoolArray
	Cars             []*Car
	ClockString      *atomic.Value
	CarsLink         *atomic.Value
	CarsTable        *CarsTable
}

func NewGlobals() *Globals {
	g := &Globals{
		InputCheckbox:    &atomic.Value{},
		InputRadioGroup1: jaws.NewNamedBoolArray().Add("1", "Radio 1.1").Add("2", "Radio 1.2"),
		InputRadioGroup2: jaws.NewNamedBoolArray().Add("1", "Radio 2.1").Add("2", "Radio 2.2"),
		InputDate:        &atomic.Value{},
		InputRange:       &atomic.Value{},
		InputButton:      &atomic.Value{},
		SelectPet:        newUiSelectPet(),
		ClockString:      &atomic.Value{},
		CarsLink:         &atomic.Value{},
		Cars: []*Car{
			{
				VIN:   "JH4DB1671PS002584",
				Year:  1993,
				Make:  "Acura",
				Model: "Integra",
			},
			{
				VIN:   "KM8JT3AC2DU583865",
				Year:  2013,
				Make:  "Hyundai",
				Model: "Tucson",
			},
			{
				VIN:   "1D4GP24R75B188657",
				Year:  2005,
				Make:  "Dodge",
				Model: "Grand Caravan",
			},
		},
	}
	g.inputTextArea = "The quick brown fox jumps over the lazy dog"
	g.InputCheckbox.Store(false)
	g.InputDate.Store(time.Now())
	g.InputRange.Store(0.0)
	g.InputButton.Store("Meh")
	g.ClockString.Store("")
	g.CarsLink.Store("...")
	return g
}

func (g *Globals) SetInputButtonID() string {
	return "setinputbutton"
}

func (g *Globals) OnSetInputButton() jaws.EventFn {
	return func(rq *jaws.Request, evt what.What, id, val string) error {
		if evt == what.Trigger {
			g.mu.RLock()
			isWoo := strings.HasPrefix(g.InputButton.Load().(string), "Woo")
			g.mu.RUnlock()
			if isWoo {
				rq.RemoveAttr(g.InputButton, "disabled")
			} else {
				rq.SetAttr(g.InputButton, "disabled", "")
			}
			g.mu.Lock()
			defer g.mu.Unlock()
			g.InputButton.Store(val)
			rq.Jaws.Dirty(g.InputButton)
		}
		return nil
	}
}

func (g *Globals) InputRangeID() string {
	return "inputrange"
}
