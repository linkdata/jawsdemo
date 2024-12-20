package main

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type Globals struct {
	FaviconURI        string
	mu                deadlock.RWMutex
	inputText         string
	inputTextArea     string
	inputCheckbox     bool
	InputRadioGroup1  *jaws.NamedBoolArray
	InputRadioGroup2  *jaws.NamedBoolArray
	inputDate         time.Time
	inputRange        float64
	inputButton       string
	SelectPet         *jaws.NamedBoolArray
	Cars              []*Car
	carsLink          string
	CarsTable         *CarsTable
	client            map[*jaws.Session]*Client
	runtime           string
	getUserAgentParam string
	userAgent         string
}

func NewGlobals() *Globals {
	g := &Globals{
		InputRadioGroup1: jaws.NewNamedBoolArray().Add("1", "Radio 1.1").Add("2", "Radio 1.2"),
		InputRadioGroup2: jaws.NewNamedBoolArray().Add("1", "Radio 2.1").Add("2", "Radio 2.2"),
		inputDate:        time.Now(),
		inputButton:      "meh",
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
		client: make(map[*jaws.Session]*Client),
	}
	g.inputTextArea = "The quick brown fox jumps over the lazy dog"
	return g
}

var _ jaws.ClickHandler = (*Globals)(nil)

func (g *Globals) JawsClick(e *jaws.Element, name string) error {
	if name == "clicky" {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.inputButton == "Meh" {
			g.inputButton = "Woo"
			e.Jaws.SetAttr(g.InputButton(), "disabled", "")
		} else {
			g.inputButton = "Meh"
			e.Jaws.RemoveAttr(g.InputButton(), "disabled")
		}
		e.Dirty(g.InputButton())
		return nil
	}
	return jaws.ErrEventUnhandled
}

func (g *Globals) Clock() jaws.HTMLGetter {
	return uiClock{}
}

func (g *Globals) Runtime() any {
	return jaws.NewJsVar(jaws.Bind(&g.mu, &g.runtime))
}

func (g *Globals) GetUserAgent() jaws.JsVar[string] {
	return jaws.NewJsVar(jaws.Bind(&g.mu, &g.getUserAgentParam))
}

func (g *Globals) UserAgent() jaws.JsVar[string] {
	return jaws.NewJsVar(jaws.Bind(&g.mu, &g.userAgent))
}

func (g *Globals) ClientUserAgent() any {
	return jaws.Bind(&g.mu, &g.userAgent)
}

func (g *Globals) JawsVersion() (v string) {
	if bi, ok := debug.ReadBuildInfo(); ok {
		v = bi.Main.Version
		for _, dep := range bi.Deps {
			if strings.HasSuffix(dep.Path, "/jaws") {
				v += " - jaws@" + dep.Version
			}
		}
	}
	return
}
