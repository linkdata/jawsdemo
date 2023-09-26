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

func (uis *UiState) Clock() jaws.HtmlGetter {
	return uiClock{}
}

func (uis *UiState) JawsVersion() (v string) {
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
	return func(rq *jaws.Request, jid string) error {
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
