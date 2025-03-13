package pipeline

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"vibe-up/sensor-data/models"
)

// ProcessFile reads and processes a single file, converting each record into DataFrameRow structs and sending them to the router
func ProcessFile(factory func([]string) (models.Sensor, error), filename string, routerChan chan<- *models.Sensor, summaryChan chan<- *models.Sensor) error {
	fmt.Println("Processing file", filename)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader for %s: %v", filename, err)
	}
	defer gzReader.Close()

	csvReader := csv.NewReader(gzReader)
	if _, err := csvReader.Read(); err != nil { // Read and discard the header
		return fmt.Errorf("failed to read header from %s: %v", filename, err)
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed to read record from %s: %v", filename, err)
			continue
		}

		record, err := factory(row)
		if err != nil {
			log.Printf("Failed to create record: %v", err)
			continue
		}

		if record.Length() > 0 {
			routerChan <- &record
		}
		summaryChan <- &record
	}

	return nil
}
