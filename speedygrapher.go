package speedygrapher

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// This package provides tools necessary to build a Graph
// using a CSV as an adjacency table.
// Below is a sample usage of the processor, but this can be
// repurposed to build graphs quickly from other data sources too.

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	finalGraph, err := CsvToGraph("sample.csv", CsvProduceRecord)
	if err != nil {
		log.Fatal(err)
	}
	err = CreateDOTFile("output_graph", finalGraph)
	if err != nil {
		log.Fatal(err)
	}
}