package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/ui"
)

func useTestGlobals(t *testing.T) *Globals {
	t.Helper()
	oldGlobals := globals
	g := NewGlobals()
	globals = g
	t.Cleanup(func() {
		globals = oldGlobals
	})
	return g
}

func newTestElement(t *testing.T) *jaws.Element {
	t.Helper()
	jw, err := jaws.New()
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	jw.NewSession(nil, req)
	rq := jw.NewRequest(req)
	t.Cleanup(func() {
		jw.Close()
	})
	return rq.NewElement(ui.NewButton("test"))
}

func TestNewGlobalsInitializesSharedUIState(t *testing.T) {
	g := NewGlobals()
	if g.CarsTable == nil {
		t.Fatal("CarsTable is nil")
	}
	if g.CarsTable.owner() != g {
		t.Fatal("CarsTable is not bound to its Globals")
	}
	if got := g.CarsTable.JawsGetTag(nil); got != &g.Cars {
		t.Fatalf("CarsTable tag = %p, want %p", got, &g.Cars)
	}
	if g.inputButton != "Meh" {
		t.Fatalf("inputButton = %q, want %q", g.inputButton, "Meh")
	}
}

func TestAddRandomCarUsesReceiverState(t *testing.T) {
	g := NewGlobals()
	before := len(g.Cars)
	car := g.AddRandomCar()
	if car == nil {
		t.Fatal("AddRandomCar returned nil")
	}
	if len(g.Cars) != before+1 {
		t.Fatalf("len(Cars) = %d, want %d", len(g.Cars), before+1)
	}
	if len(car.VIN) != 17 {
		t.Fatalf("VIN length = %d, want 17", len(car.VIN))
	}
}

func TestCarClickReordersAndRemovesWithGlobalLock(t *testing.T) {
	g := useTestGlobals(t)
	elem := newTestElement(t)
	first, second, third := g.Cars[0], g.Cars[1], g.Cars[2]

	if err := second.JawsClick(elem, jaws.Click{Name: "up"}); err != nil {
		t.Fatal(err)
	}
	if got := g.Cars; got[0] != second || got[1] != first || got[2] != third {
		t.Fatalf("after up = [%p %p %p], want [%p %p %p]", got[0], got[1], got[2], second, first, third)
	}

	if err := second.JawsClick(elem, jaws.Click{Name: "down"}); err != nil {
		t.Fatal(err)
	}
	if got := g.Cars; got[0] != first || got[1] != second || got[2] != third {
		t.Fatalf("after down = [%p %p %p], want [%p %p %p]", got[0], got[1], got[2], first, second, third)
	}

	if err := second.JawsClick(elem, jaws.Click{Name: "remove"}); err != nil {
		t.Fatal(err)
	}
	if got := g.Cars; len(got) != 2 || got[0] != first || got[1] != third {
		t.Fatalf("after remove = %#v, want first and third cars", got)
	}
}

func TestCarClickConditionBounds(t *testing.T) {
	elem := newTestElement(t)
	car := &Car{condition: 99}

	if err := car.JawsClick(elem, jaws.Click{Name: "+"}); err != nil {
		t.Fatal(err)
	}
	if car.condition != 100 {
		t.Fatalf("condition = %v, want 100", car.condition)
	}
	if err := car.JawsClick(elem, jaws.Click{Name: "+"}); err == nil {
		t.Fatal("expected high condition error")
	}

	car.condition = 1
	if err := car.JawsClick(elem, jaws.Click{Name: "-"}); err != nil {
		t.Fatal(err)
	}
	if car.condition != 0 {
		t.Fatalf("condition = %v, want 0", car.condition)
	}
	if err := car.JawsClick(elem, jaws.Click{Name: "-"}); err == nil {
		t.Fatal("expected low condition error")
	}
}

func TestInputButtonFirstClickEntersMysticalState(t *testing.T) {
	g := useTestGlobals(t)
	elem := newTestElement(t)

	if err := g.InputButton().JawsClick(elem, jaws.Click{Name: "clicky"}); err != nil {
		t.Fatal(err)
	}
	if g.inputButton != "Mystical?" {
		t.Fatalf("inputButton = %q, want %q", g.inputButton, "Mystical?")
	}
	if got := elem.Session().Get("mystical"); got != true {
		t.Fatalf("mystical session value = %#v, want true", got)
	}

	if err := g.CarsTable.Mystical().JawsClick(elem, jaws.Click{}); err != nil {
		t.Fatal(err)
	}
	if got := elem.Session().Get("mystical"); got != nil {
		t.Fatalf("mystical session value = %#v, want nil", got)
	}
}

func TestClientJsVarDoesNotRequireInitialSession(t *testing.T) {
	g := useTestGlobals(t)
	jw, err := jaws.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(jw.Close)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rq := jw.NewRequest(req)
	jsvar, err := g.Client().JawsMakeJsVar(rq)
	if err != nil {
		t.Fatal(err)
	}
	if rq.Session() != nil {
		t.Fatal("initial render created a session")
	}
	clientVar, ok := jsvar.(*ui.JsVar[Client])
	if !ok {
		t.Fatalf("JawsMakeJsVar returned %T, want *ui.JsVar[Client]", jsvar)
	}
	if got := clientVar.JawsGet(nil); got != (Client{X: -1, Y: -1}) {
		t.Fatalf("initial client = %#v, want zero value", got)
	}
}

func TestClientPathSetStoresClientInSession(t *testing.T) {
	elem := newTestElement(t)
	client := &Client{X: 12, Y: 34, B: 1}

	client.JawsPathSet(elem, "X", client.X)
	if got := elem.Session().Get(clientSessionKey); got != client {
		t.Fatalf("client session value = %#v, want %p", got, client)
	}
}

func TestRoutesRender(t *testing.T) {
	g := useTestGlobals(t)
	jw, err := jaws.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(jw.Close)
	jw.AutoSession = true

	mux := http.NewServeMux()
	if err := setupRoutes(jw, mux); err != nil {
		t.Fatal(err)
	}
	g.FaviconURI = jw.FaviconURL()
	handler := jw.SecureHeadersMiddleware(mux)

	for _, path := range []string{"/", "/cars"} {
		t.Run(path, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", rr.Code, http.StatusOK, rr.Body.String())
			}
			if !strings.Contains(rr.Body.String(), "<!doctype html>") {
				t.Fatalf("response for %q did not render an HTML document", path)
			}
			if cookies := rr.Result().Cookies(); len(cookies) != 0 {
				t.Fatalf("response for %q created cookies: %#v", path, cookies)
			}
		})
	}
}
