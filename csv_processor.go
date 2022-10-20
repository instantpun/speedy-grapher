package speedygrapher

import (
	"encoding/csv"
	"fmt"
	"log"
	"sync"
	"io"
	"github.com/awalterschulze/gographviz"
)

func CsvToGraph(infile, outfilePrefix string, recordProducer func(*io.Reader,*FieldnameMap) (Record, error)) (*gographviz.Interface, error) {

	fh, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	reader := csv.NewReader(fh)
	fnames, err := ExtractFieldnames(reader)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	graphOrch := NewGraphOrchestrator()

	i := 0
	for {
		record, err := recordProducer(reader, fnames)
		if err == io.EOF {
			log.Println("Reached end of source file")
			break
		} else if err != nil {
			return nil, err
		} else {
			log.Debug("Record %d: %v\n", i, record)

			// Deadlock prevention:
			// goroutines will remain alive, but asleep indefinitely until
			// something reads from the unbuffered error channel 
			go func() {
				for asyncErr := range graphOrch.Errors {
						if asyncErr != nil {
							return nil, err // return first error from go routine
						}
					}
			}()

			CoordinatedGraphUpdate(record, graphOrch, &wg)
		}
		i++

	}

	return graphOrch.Graph
}

func CreateDOTFile(g *gographviz.Interface) error {
	log.Println("Writing to file...")
	log.Debug(g.String())
	
	outputGraphFile, err := os.Create(fmt.Sprintf("%s.dot", outfilePrefix))
	if err != nil {
		return err
	}
	defer outputGraphFile.Close()
	
	outputGraphFile.Write([]byte(g.String()))
	
	return nil
}


func CsvExtractFieldnames(fh *io.Reader) (*FieldnameMap, error) {

	fnames := NewFieldnameMap()
	
	firstLine, err := fh.Read()
	if err != nil {
		return nil, err // any I/O error or io.EOF signal
	}
	
	fnames, err = fnames.UpdateFromSlice(firstLine)
	if err != nil {
		return nil, err
	}

	return fnames, nil
}

func CsvProduceRecord(fh *csv.Reader, fnames *FieldnameMap) (Record, error) {
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

