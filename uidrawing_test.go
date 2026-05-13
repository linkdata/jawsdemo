package main

import (
	"math"
	"strings"
	"testing"

	"github.com/linkdata/jaws"
)

func TestDrawingPointerDownAddsSVGPoint(t *testing.T) {
	g := NewGlobals()
	elem := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	drawTestPoint(t, drawing, elem, 12.25, 34.75)
	if len(g.drawingTraces) != 1 {
		t.Fatalf("drawing traces = %d, want 1", len(g.drawingTraces))
	}
	trace := g.drawingTraces[0]
	if trace.SessionID != elem.Session().ID() {
		t.Fatalf("session ID = %v, want %v", trace.SessionID, elem.Session().ID())
	}
	if trace.Color != drawingColors[0] {
		t.Fatalf("color = %q, want %q", trace.Color, drawingColors[0])
	}
	if len(trace.Points) != 1 {
		t.Fatalf("trace points = %d, want 1", len(trace.Points))
	}
	if got := trace.Points[0]; got.X != 12.25 || got.Y != 34.75 {
		t.Fatalf("point = %+v, want 12.25,34.75", got)
	}

	html := string(drawing.JawsGetHTML(elem))
	if !strings.Contains(html, `<circle cx="12.2" cy="34.8"`) {
		t.Fatalf("drawing HTML does not contain rounded point: %s", html)
	}
	if !strings.Contains(html, `fill="`+drawingColors[0]+`"`) {
		t.Fatalf("drawing HTML does not contain trace color %q: %s", drawingColors[0], html)
	}
}

func TestDrawingPointerMoveAddsSVGPointWhilePrimaryButtonHeld(t *testing.T) {
	g := NewGlobals()
	elem := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	drawTestPoint(t, drawing, elem, 1, 2)
	if err := drawing.JawsPointer(elem, jaws.Pointer{
		Name:    "drawing",
		X:       3,
		Y:       4,
		Kind:    jaws.PointerMove,
		Button:  -1,
		Buttons: jaws.PointerButtonPrimary,
	}); err != nil {
		t.Fatal(err)
	}
	trace := findTestTrace(t, g, elem.Session().ID())
	if len(trace.Points) != 2 {
		t.Fatalf("trace points = %d, want 2", len(trace.Points))
	}
	if got := trace.Points[1]; got.X != 3 || got.Y != 4 {
		t.Fatalf("move point = %+v, want 3,4", got)
	}
}

func TestDrawingPointerMoveWithoutButtonIsIgnored(t *testing.T) {
	g := NewGlobals()
	elem := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	if err := drawing.JawsPointer(elem, jaws.Pointer{
		Name:   "drawing",
		X:      3,
		Y:      4,
		Kind:   jaws.PointerMove,
		Button: -1,
	}); err != nil {
		t.Fatal(err)
	}
	if len(g.drawingTraces) != 0 {
		t.Fatalf("drawing traces = %d, want 0", len(g.drawingTraces))
	}
}

func TestDrawingTracksSeparateSessionTraces(t *testing.T) {
	g := NewGlobals()
	elem1 := newTestElement(t)
	elem2 := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	drawTestPoint(t, drawing, elem1, 1, 2)
	drawTestPoint(t, drawing, elem2, 3, 4)

	trace1 := findTestTrace(t, g, elem1.Session().ID())
	trace2 := findTestTrace(t, g, elem2.Session().ID())
	if trace1 == trace2 {
		t.Fatal("expected separate trace pointers")
	}
	if trace1.Color == trace2.Color {
		t.Fatalf("expected distinct colors, both were %q", trace1.Color)
	}
	if got := trace1.Points[0]; got.X != 1 || got.Y != 2 {
		t.Fatalf("trace1 point = %+v, want 1,2", got)
	}
	if got := trace2.Points[0]; got.X != 3 || got.Y != 4 {
		t.Fatalf("trace2 point = %+v, want 3,4", got)
	}

	html := string(drawing.JawsGetHTML(elem1))
	if !strings.Contains(html, `stroke="`+trace1.Color+`"`) || !strings.Contains(html, `stroke="`+trace2.Color+`"`) {
		t.Fatalf("drawing HTML does not contain both trace colors: %s", html)
	}
}

func TestDrawingIgnoresOtherPointerNames(t *testing.T) {
	g := NewGlobals()
	elem := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	if err := drawing.JawsPointer(elem, jaws.Pointer{Name: "other", X: 12, Y: 34, Kind: jaws.PointerDown, Button: 0}); err != jaws.ErrEventUnhandled {
		t.Fatalf("error = %v, want ErrEventUnhandled", err)
	}
	if len(g.drawingTraces) != 0 {
		t.Fatalf("drawing traces = %d, want 0", len(g.drawingTraces))
	}
}

func TestDrawingContextMenuClearsCurrentSessionTrace(t *testing.T) {
	g := NewGlobals()
	elem1 := newTestElement(t)
	elem2 := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	drawTestPoint(t, drawing, elem1, 1, 2)
	drawTestPoint(t, drawing, elem2, 3, 4)
	if err := drawing.JawsContextMenu(elem1, jaws.Click{Name: "drawing"}); err != nil {
		t.Fatal(err)
	}

	trace1 := findTestTrace(t, g, elem1.Session().ID())
	trace2 := findTestTrace(t, g, elem2.Session().ID())
	if len(trace1.Points) != 0 {
		t.Fatalf("trace1 points = %d, want 0", len(trace1.Points))
	}
	if len(trace2.Points) != 1 {
		t.Fatalf("trace2 points = %d, want 1", len(trace2.Points))
	}
}

func TestDrawingClampsCoordinates(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		max  float64
		want float64
	}{
		{name: "negative", in: -1, max: drawingWidth, want: 0},
		{name: "nan", in: math.NaN(), max: drawingWidth, want: 0},
		{name: "inf", in: math.Inf(1), max: drawingWidth, want: 0},
		{name: "over", in: drawingWidth + 1, max: drawingWidth, want: drawingWidth},
		{name: "inside", in: 42, max: drawingWidth, want: 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clampDrawingCoord(tt.in, tt.max); got != tt.want {
				t.Fatalf("clampDrawingCoord(%v, %v) = %v, want %v", tt.in, tt.max, got, tt.want)
			}
		})
	}
}

func TestDrawingRetainsMaximumPoints(t *testing.T) {
	g := NewGlobals()
	elem := newTestElement(t)
	drawing := g.Drawing().(uiDrawing)

	for i := 0; i < drawingMaxPoints+3; i++ {
		drawTestPoint(t, drawing, elem, float64(i), float64(i))
	}
	trace := findTestTrace(t, g, elem.Session().ID())
	if len(trace.Points) != drawingMaxPoints {
		t.Fatalf("drawing points = %d, want %d", len(trace.Points), drawingMaxPoints)
	}
	if got := trace.Points[0]; got.X != 3 || got.Y != 3 {
		t.Fatalf("first retained point = %+v, want 3,3", got)
	}
}

func drawTestPoint(t *testing.T, drawing uiDrawing, elem *jaws.Element, x, y float64) {
	t.Helper()
	if err := drawing.JawsPointer(elem, jaws.Pointer{
		Name:    "drawing",
		X:       x,
		Y:       y,
		Kind:    jaws.PointerDown,
		Button:  0,
		Buttons: jaws.PointerButtonPrimary,
	}); err != nil {
		t.Fatal(err)
	}
}

func findTestTrace(t *testing.T, g *Globals, sessionID uint64) *drawingTrace {
	t.Helper()
	for i := range g.drawingTraces {
		if g.drawingTraces[i].SessionID == sessionID {
			return &g.drawingTraces[i]
		}
	}
	t.Fatalf("missing trace for session %v", sessionID)
	return nil
}
