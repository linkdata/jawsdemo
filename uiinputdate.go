package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type uiInputDate struct {
	jid  string
	mu   deadlock.RWMutex // protects following
	data time.Time
}

func newUiInputDate(jid string) *uiInputDate {
	return &uiInputDate{jid: jid}
}

// eventFn gets called by JaWS when the client browser Javascript reports that the data has changed.
func (ui *uiInputDate) eventFn(rq *jaws.Request, jid string, val time.Time) error {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	// it's usually a good idea to ensure that the value is actually changed before doing work
	if ui.data != val {
		ui.data = val
		// sends the changed text to all *other* Requests.
		rq.SetDateValue(ui.jid, val)
	}
	return nil
}

func (ui *uiInputDate) JawsUi(rq *jaws.Request, attrs ...string) template.HTML {
	ui.mu.RLock()
	data := ui.data
	ui.mu.RUnlock()
	return template.HTML(fmt.Sprintf(`<label for="%s" class="form-label">Date</label>`, ui.jid)) +
		rq.Date(ui.jid, data, []interface{}{ui.eventFn, attrs})
}
