package pipeline

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"

	"vibe-up/sensor-data/models"
)

func SummaryStatistics(summaryChan chan *models.Sensor) {

	summaryFile, err := os.Create("output/summary_output.csv")

	if err != nil {
		log.Fatalf("Failed to create summary file: %v", err)
	}
	defer summaryFile.Close()

	bufferedWriter := bufio.NewWriter(summaryFile)
	csvWriter := csv.NewWriter(bufferedWriter)

	if err := csvWriter.Write([]string{"SubmissionID", "ParticipantID", "TrialID", "Timestamp", "Timezone", "N", "MaxOffset", "RMS"}); err != nil {
		log.Fatalf("Failed to write header to summary file XXX: %v", err)
	}

	defer csvWriter.Flush()
	defer bufferedWriter.Flush()

	for record := range summaryChan {
		summary := (*record).Summarise()

		if err := csvWriter.Write(summary); err != nil {
			log.Printf("Failed to write summary: %v", err)
			continue
		}
	}
	log.Println("Summary statistics written to output/summary_output.csv")
}
