package main

import (
	"html/template"

	"github.com/linkdata/jaws"
)

type uiSelectPet struct {
	nba *jaws.NamedBoolArray
}

func newUiSelectPet(jid string) *uiSelectPet {
	nba := jaws.NewNamedBoolArray(jid)
	nba.Add("", "--Please choose an option--")
	nba.Add("dog", "Dog")
	nba.Add("cat", "Cat")
	nba.Add("hamster", "Hamster")
	nba.Add("parrot", "Parrot")
	nba.Add("spider", "Spider")
	return &uiSelectPet{
		nba: nba,
	}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiSelectPet) eventFn(rq *jaws.Request, jid, val string) error {
	if val != ui.nba.Get() {
		ui.nba.SetOnly(val)
		rq.SetTextValue(ui.nba.Jid, val)
	}
	return nil
}

func (ui *uiSelectPet) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	return rq.Select(ui.nba, ui.eventFn, attrs...)
}
