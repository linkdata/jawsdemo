package main

import (
	"runtime/debug"
	"strings"

	"github.com/linkdata/jaws"
)

type UiState struct {
	G *Globals
	*jaws.Request
}

func (uis *UiState) UiClock() jaws.Ui         { return &uiClock{uis.G} }
func (uis *UiState) UiCarsLink() jaws.Ui      { return &uiCarsLink{uis.G} }
func (uis *UiState) UiInputText() jaws.Ui     { return &uiInputText{uis.G} }
func (uis *UiState) UiInputRange() jaws.Ui    { return &uiInputRange{uis.G} }
func (uis *UiState) UiInputCheckbox() jaws.Ui { return &uiInputCheckbox{uis.G} }
func (uis *UiState) UiInputDate() jaws.Ui     { return &uiInputDate{uis.G} }

func (uis *UiState) Version() (v string) {
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

func (uis *UiState) OnMystical() jaws.ClickFn {
	return func(rq *jaws.Request) error {
		rq.Trigger("setinputbutton", "Wooooo....")
		return nil
	}
}

func NewUiState(rq *jaws.Request, g *Globals) *UiState {
	uis := &UiState{
		G:       g,
		Request: rq,
	}
	return uis
}
