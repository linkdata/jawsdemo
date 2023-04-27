package main

import (
	"log"

	"github.com/linkdata/jaws"
)

type uiInputRadioGroup struct {
	nba *jaws.NamedBoolArray
}

func newUiInputRadioGroup(nba *jaws.NamedBoolArray) *uiInputRadioGroup {
	return &uiInputRadioGroup{
		nba: nba,
	}
}

func (ui *uiInputRadioGroup) JawsRadioGroupData() *jaws.NamedBoolArray {
	return ui.nba
}

func (ui *uiInputRadioGroup) JawsRadioGroupHandler(rq *jaws.Request, boolName string) error {
	log.Println("uiInputRadioGroup:", ui.nba)
	return nil
}
