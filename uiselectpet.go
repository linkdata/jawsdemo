package main

import (
	"fmt"
	"html/template"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiSelectPet struct {
	jid  string
	mu   deadlock.RWMutex // protects following
	data string
}

func newUiSelectPet(jid string) *uiSelectPet {
	return &uiSelectPet{
		jid: jid,
	}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiSelectPet) eventFn(rq *jaws.Request, val string) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if ui.data != val {
		ui.data = val
		rq.SetTextValue(ui.jid, val)
	}
	return nil
}

func (ui *uiSelectPet) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	ui.mu.RLock()
	data := ui.data
	ui.mu.RUnlock()
	nba := jaws.NewNamedBoolArray()
	nba.Add("", "--Please choose an option--")
	nba.Add("dog", "Dog")
	nba.Add("cat", "Cat")
	nba.Add("hamster", "Hamster")
	nba.Add("parrot", "Parrot")
	nba.Add("spider", "Spider")
	nba.Check(data)
	return template.HTML(fmt.Sprintf(`<label for="%s">Choose a pet:</label>%s`,
		ui.jid,
		rq.Select(ui.jid, nba, ui.eventFn, attrs...),
	))
}
