package main

import (
	"html/template"
	"runtime/debug"
	"strings"

	"github.com/linkdata/jaws"
)

type UiState struct {
	G *Globals
	*jaws.Request
}

func (uis *UiState) GetHeadFooter() template.HTML {
	return uis.HeadHTML()
}

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
