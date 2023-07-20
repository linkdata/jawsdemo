package main

import (
	"log"

	"github.com/linkdata/jaws"
)

type uiInputRadioGroup struct {
	*jaws.NamedBoolArray
}

func newUiInputRadioGroup(nba *jaws.NamedBoolArray) *uiInputRadioGroup {
	return &uiInputRadioGroup{nba}
}

// JawsRadioGroupData has a default implementation in jaws.NamedBoolArray
// and so does JawsRadioGroupHandler, but we override it here to print the state.
func (ui *uiInputRadioGroup) JawsRadioGroupHandler(rq *jaws.Request, jid, boolName string) error {
	log.Println("uiInputRadioGroup:", ui.NamedBoolArray)
	return nil
}
