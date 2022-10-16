package main

import (
	"encoding/csv"
	"fmt"
	"sync"
	"github.com/awalterschulze/gographviz"
)

type GraphAsync struct {
	gographviz.Graph
	sync.RWMutex
}

func ExtractFieldnames(fh *csv.Reader) (*FieldnameMap, error) {
	fnames := NewFieldnameMap()
	
	firstLine, err := fh.Read()
	if err != nil {
		return nil, err // any I/O error or io.EOF signal
	}
	
	// var ft interface{} = firstLine
	// _, ok := ft.([]string)
	// if !ok {
	// 	return nil, fmt.Errorf("Error during processing. Expected []string found %T", firstLine)
	// }

	fnames, err = fnames.Update(firstLine)
	if err != nil {
		return nil, err
	}

	return fnames, nil
}

func ProduceCsvRecord(fh *csv.Reader, fnames *FieldnameMap) (Record, error) {
	if fh == nil {
		return nil, fmt.Errorf("Error Producing Record: *csv.Reader is nil pointer.")
	}

	if fnames == nil {
		return nil, fmt.Errorf("Error Producing Record: *FielnameMap is nil pointer.")
	}

	rowValues, err := fh.Read()
	if err != nil {
		return nil, err // any I/O error or io.EOF signal
	}

	r := Record{}
	for i, val := range rowValues {
		r[(*fnames)[i]] = val
	}
	return r, nil
}

func UpdateGraphByRecord(
	record Record, 
	graph *GraphAsync, 
	parentGraphName string, 
	srcSet, dstSet *StringSetAsync, 
	edgeSet *EdgeSetAsync,
	wg *sync.WaitGroup) (error) {
	
	wg.Add(1)
	go func(record Record, parentGraphName string) error {
		defer wg.Done()

		if !srcSet.Has(record["source"]) {
			srcSet.Add(record["source"])
		
			graph.Lock()
			defer graph.Unlock()
			err := graph.AddNode(parentGraphName, record["source"], nil)
			if err != nil {
				return err
			}
		}		
		return nil
	}(record, parentGraphName)

	wg.Add(1)
	go func(record Record, parentGraphName string) error {
		defer wg.Done()

		if !dstSet.Has(record["destination"]) {
			dstSet.Add(record["destination"])

			graph.Lock()
			defer graph.Unlock()
			err := graph.AddNode(parentGraphName, record["destination"], nil)
			if err != nil {
				return err
			}
		}
	
		return nil
	}(record, parentGraphName)

	wg.Add(1)
	go func(record Record) error {
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
				return err
			}
		}

		return nil
	}(record)

	wg.Wait()

	return nil
}