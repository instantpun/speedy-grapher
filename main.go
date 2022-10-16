package main

import (
	// "fmt"
	"os"
	"encoding/csv"
	"io"
	"github.com/awalterschulze/gographviz"
	"log"
	"sync"
)

func main() {
	fh, err := os.Open("sample.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	reader := csv.NewReader(fh)

	// seutp graph
	graphAst, _ := gographviz.ParseString(`digraph G { rankdir=LR;}`)
	graph := gographviz.NewGraph()
	err = gographviz.Analyse(graphAst, graph)
	if err != nil {
		log.Fatal(err)
	}
	lockedGraph := &GraphAsync{
		*graph,
		sync.RWMutex{},
	}

	fnames, err := ExtractFieldnames(reader)
	if err != nil {
		log.Fatal(err)
	}

	sources := NewStringSetAsync()
	destinations := NewStringSetAsync()
	edges := NewEdgeSetAsync()

	wg := sync.WaitGroup{}

	i := 0
	for {
		// process records
		record, err := ProduceCsvRecord(reader, fnames)
		if err == io.EOF {
			log.Println("End of document")
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Record %d: %v\n", i, record)
		}

		// 
		err = UpdateGraphByRecord(
			record, lockedGraph, "G",
			sources, destinations, edges, &wg)
		if err != nil {
			log.Fatal(err)
		} 
		i++
	}

	log.Println(graph.String())
}