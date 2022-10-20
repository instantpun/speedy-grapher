package speedygrapher

import (
	"encoding/csv"
	"fmt"
	"sync"
	"io"
	"os"
	"github.com/awalterschulze/gographviz"
)

func CsvToGraph(infile string, recordProducer func(*csv.Reader,*FieldnameMap) (Record, error)) (gographviz.Interface, error) {

	fh, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	reader := csv.NewReader(fh)
	fnames, err := CsvExtractFieldnames(reader)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	graphOrch, err := NewGraphOrchestrator()
	if err != nil {
		return nil, err
	}

	i := 0
	for {
		record, err := recordProducer(reader, fnames)
		if err == io.EOF {
			fmt.Printf("Reached end of source file. Total Records: %d\n", i)
			break
		} else if err != nil {
			return nil, err
		} else {
			// Deadlock prevention:
			// goroutines will remain alive, but asleep indefinitely until
			// something reads from the unbuffered error channel 
			go func() {
				for asyncErr := range graphOrch.Errors {
						if asyncErr != nil {
							panic(err)
						}
					}
			}()

			CoordinatedGraphUpdate(record, graphOrch, &wg)
		}
		i++

	}

	return graphOrch.Graph, nil
}

func CreateDOTFile(filePrefix string, g gographviz.Interface) error {
	fmt.Println("Writing to file...")
	
	outputGraphFile, err := os.Create(fmt.Sprintf("%s.dot", filePrefix))
	if err != nil {
		return err
	}
	defer outputGraphFile.Close()
	
	outputGraphFile.Write([]byte(g.String()))
	
	return nil
}


func CsvExtractFieldnames(fh *csv.Reader) (*FieldnameMap, error) {

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

