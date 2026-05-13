package main

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws/lib/bind"
	"github.com/linkdata/jaws/lib/named"
	"github.com/linkdata/jaws/lib/ui"
)

type Globals struct {
	FaviconURI       string
	mu               deadlock.RWMutex
	inputText        string
	inputTextArea    string
	inputCheckbox    bool
	InputRadioGroup1 *named.BoolArray
	InputRadioGroup2 *named.BoolArray
	inputDate        time.Time
	inputRange       float64
	inputButton      string
	SelectPet        *named.BoolArray
	Cars             []*Car
	carsLink         string
	CarsTable        *CarsTable
	drawingTraces    []drawingTrace
	runtime          string
}

func NewGlobals() *Globals {
	g := &Globals{
		InputRadioGroup1: named.NewBoolArray(false).Add("1", "Radio 1.1").Add("2", "Radio 1.2"),
		InputRadioGroup2: named.NewBoolArray(false).Add("1", "Radio 2.1").Add("2", "Radio 2.2"),
		inputDate:        time.Now(),
		inputButton:      "Meh",
		SelectPet:        newUiSelectPet(),
		carsLink:         "...",
		Cars: []*Car{
			{
				VIN:       "JH4DB1671PS002584",
				Year:      1993,
				Make:      "Acura",
				Model:     "Integra",
				condition: 12,
			},
			{
				VIN:       "KM8JT3AC2DU583865",
				Year:      2013,
				Make:      "Hyundai",
				Model:     "Tucson",
				condition: 97,
			},
			{
				VIN:       "1D4GP24R75B188657",
				Year:      2005,
				Make:      "Dodge",
				Model:     "Grand Caravan",
				condition: 67,
			},
		},
	}
	g.CarsTable = &CarsTable{globals: g}
	g.inputTextArea = "The quick brown fox jumps over the lazy dog"
	return g
}

func (g *Globals) Clock() bind.HTMLGetter {
	return uiClock{}
}

func (g *Globals) Runtime() any {
	return ui.NewJsVar(&g.mu, &g.runtime)
}

func (g *Globals) JawsVersion() (v string) {
	v = PkgVersion
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range bi.Deps {
			if strings.HasSuffix(dep.Path, "/jaws") {
				v += " - jaws@" + dep.Version
			}
		}
	}
	return
}
