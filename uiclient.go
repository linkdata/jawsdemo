package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/ui"
)

type clientState struct {
	value Client
	store *ui.JsVarStore[Client]
}

// ClientBinding renders the request's browser client state.
func (g *Globals) ClientBinding(rq *jaws.Request) (binding *ui.JsVarBinding[Client], err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	state, _ := rq.Get(clientSessionKey).(*clientState)
	if state == nil {
		state = &clientState{value: Client{X: -1, Y: -1}}
		if state.store, err = ui.NewJsVarStore(rq.Jaws, "client", &g.mu, &state.value); err != nil {
			return
		}
		state.store.ClientCheck = func(*jaws.Element, *Client, string) error { return nil }
		state.store.ExtraTags = []any{uiClientPos{}}
		if rq.Session() != nil {
			rq.Set(clientSessionKey, state)
		} else {
			rq.SetConnectFn(func(rq *jaws.Request) error {
				rq.Set(clientSessionKey, state)
				return nil
			})
		}
	}
	binding = state.store.Bind()
	return
}
