package pipeline

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"vibe-up/sensor-data/models"
)

// CSVWriter receives records from a sink and prints the summary statistics
func CSVWriter(id string, sinkChan <-chan *models.Sensor, wg *sync.WaitGroup) {
	defer wg.Done()

	outputFile, err := os.Create(fmt.Sprintf("output/%s.csv.gz", id))
	if err != nil {
		log.Printf("FileWriter %s: Failed to create output file: %v", id, err)
		return
	}
	defer outputFile.Close()

	bufWriter := bufio.NewWriter(outputFile)
	defer bufWriter.Flush()

	gzipWriter := gzip.NewWriter(bufWriter)
	defer gzipWriter.Close()

	writer := csv.NewWriter(gzipWriter)
	defer writer.Flush()

	// Write CSV header
	if err := writer.Write([]string{"SubmissionID", "ParticipantID", "TrialID", "Timestamp", "Timezone", "TimeOffset", "XAxis", "YAxis", "ZAxis"}); err != nil {
		log.Printf("FileWriter %s: Failed to write header: %v", id, err)
		return
	}

	for record := range sinkChan {
		rows, err := (*record).Unnest()
		if err != nil {
			log.Printf("FileWriter %s: Failed to convert record to rows: %v", id, err)
			continue
		}

		for _, row := range rows {
			if err := writer.Write(row); err != nil {
				log.Printf("FileWriter %s: Failed to write row: %v", id, err)
				return
			}
		}
	}
}
