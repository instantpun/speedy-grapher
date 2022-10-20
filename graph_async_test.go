package speedygrapher

import (
	"testing"
)

func TestNewGraphAsync(t *testing.T) {
	g := NewGraphAsync()
	t.Log(g)
}

func TestNewGraphOrchestrator(t *testing.T) {
	o, err := NewGraphOrchestrator()
	if err != nil {
		t.Error(err)
	} else if o.Graph == nil {
		t.Error("Nil pointer @ GraphOrchestrator.Graph")
	} else {
		t.Log(o.Graph.String())
	}
}

// TODO: 
// Add coverage for other functions in graph_async