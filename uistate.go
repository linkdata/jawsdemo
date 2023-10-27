package main

import (
	"io"
	"runtime/debug"
	"strings"

	"github.com/linkdata/jaws"
)

type UiState struct {
	G *Globals
	jaws.RequestWriter
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

func NewUiState(w io.Writer, rq *jaws.Request, g *Globals) *UiState {
	uis := &UiState{
		G:             g,
		RequestWriter: rq.Writer(w),
	}
	return uis
}
