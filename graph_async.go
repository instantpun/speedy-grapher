package speedygrapher

import (
	"github.com/awalterschulze/gographviz"
	"fmt"
	"sync"
	// log "github.com/sirupsen/logrus"
)
type GraphAsync struct {
	gographviz.Graph
	sync.RWMutex
}

type GraphOrchestrator struct {
	Metadata map[string]string // only use this for arbitrary metadata not compatible with DOT attributes; otherwise use Graph.Attrs
	Graph    *GraphAsync
	SrcSet   *StringSetAsync
	DstSet   *StringSetAsync
	EdgeSet  *EdgeSetAsync
	Errors   chan error
}

func NewGraphAsync() (*GraphAsync) {
	return &GraphAsync{
		*gographviz.NewGraph(),
		sync.RWMutex{},
	}
}

func NewGraphOrchestrator() (*GraphOrchestrator, error) {

	graphAst, _ := gographviz.ParseString(`digraph G { rankdir=LR;}`)
	graphAsync := NewGraphAsync()
	// Analyse() expects a type which satisfies gographviz.Interface
	// GraphAsync satisfies this via embedding
	err := gographviz.Analyse(graphAst, graphAsync)

	if err != nil {
		return nil, fmt.Errorf("In gographviz.Analyse: %v\n",err)
	}

	o := &GraphOrchestrator{
		Graph: graphAsync,
		Metadata: map[string]string{},
		SrcSet: NewStringSetAsync(),
		DstSet: NewStringSetAsync(),
		EdgeSet: NewEdgeSetAsync(),
		Errors: make(chan error, 10),
	}
	return o, nil
}

func CoordinatedGraphUpdate(record Record, o *GraphOrchestrator, wg *sync.WaitGroup) {
	
	// NOTE:
	// This spawns 3 producer goroutines, 1 for each src, dst, and edge in the record
	// However... 
	// CAUTION: 
	// None of the producer goroutines close the GraphOrchestrator.Errors channel
	// If there is not an active receiver to read from the channel BEFORE the call to sync.WaitGroup.Wait() (below)
	// then the goroutines are blocked indefinitely e.g. deadlock
	// Enjoy! :)
	wg.Add(1)
	go AddSrcNodeAsync(record, o.Graph, o.SrcSet, o.Errors, wg)
	wg.Add(1)
	go AddDstNodeAsync(record, o.Graph, o.DstSet, o.Errors, wg)
	wg.Add(1)
	go AddEdgeAsync(record, o.Graph, o.EdgeSet, o.Errors, wg)

	// synchronize goroutines
	wg.Wait()

	return
}

func AddSrcNodeAsync(record Record, graph *GraphAsync, srcSet *StringSetAsync, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	if !srcSet.Has(record["source"]) {
		srcSet.Add(record["source"])
	
		graph.Lock()
		defer graph.Unlock()

		err := graph.AddNode(graph.Name, record["source"], nil)
		if err != nil {
			errChan <- fmt.Errorf("AddSrcNodeAsync(): %v\n",err)
		}
	}		

	errChan <- nil
}

func AddDstNodeAsync(record Record, graph *GraphAsync, dstSet *StringSetAsync, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	if !dstSet.Has(record["destination"]) {
		dstSet.Add(record["destination"])

		graph.Lock()
		defer graph.Unlock()

		err := graph.AddNode(graph.Name, record["destination"], nil)
		if err != nil {
			errChan <- fmt.Errorf("AddDstNodeAsync(): %v\n",err)
		}
	}		
	errChan <- nil
}

func AddEdgeAsync(record Record, graph *GraphAsync, edgeSet *EdgeSetAsync, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	var currentEdge = AnnotatedEdge{
		record["source"],
		record["destination"],
		record["label"],
	}

	if !edgeSet.Has(currentEdge) {
		edgeSet.Add(currentEdge)
		edgeAttrs := map[string]string{	"label": currentEdge.Label, }

		graph.Lock()
		defer graph.Unlock()

		err := graph.AddEdge(
			currentEdge.Src,
			currentEdge.Dst,
			true,
			edgeAttrs,
		)
		if err != nil {
			errChan <- fmt.Errorf("AddEdgeAsync(): %v\n",err)
		}
	}

	errChan <- nil
}