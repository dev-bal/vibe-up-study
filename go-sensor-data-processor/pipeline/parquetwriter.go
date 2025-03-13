package pipeline

import (
	"fmt"
	"log"
	"sync"
	"vibe-up/sensor-data/models"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

func ParquetWriter(id string, sinkChan <-chan models.Sensor, wg *sync.WaitGroup) {
	defer wg.Done()
	// Define the schema
	md := []string{
		"name=SubmissionID, type=BYTE_ARRAY, convertedtype=UTF8",
		"name=ParticipantID, type=BYTE_ARRAY, convertedtype=UTF8",
		"name=TrialID, type=BYTE_ARRAY, convertedtype=UTF8",
		"name=Timestamp, type=BYTE_ARRAY, convertedtype=UTF8",
		"name=Timezone, type=BYTE_ARRAY, convertedtype=UTF8",
		"name=TimeOffset, type=FLOAT",
		"name=XAxis, type=FLOAT",
		"name=YAxis, type=FLOAT",
		"name=ZAxis, type=FLOAT",
	}

	// Create a Parquet file writer
	fw, err := local.NewLocalFileWriter(fmt.Sprintf("output/%s.parquet", id))
	if err != nil {
		log.Fatalf("Error creating Parquet file writer: %v", err)
	}
	defer fw.Close()

	pw, err := writer.NewCSVWriter(md, fw, 4)
	if err != nil {
		log.Fatalf("Error creating CSVWriter: %v", err)
	}

	// Write each row to the CSVWriter
	for record := range sinkChan {
		_, err := record.Unnest()
		if err != nil {
			log.Printf("Error converting record to rows: %v", err)
			continue
		}
		testrows := [][]string{
			{"SubmissionID_1", "ParticipantID_1", "TrialID_1", "2024-11-18T12:00:00Z", "+10:00", "0.1", "0.2", "0.3"},
			{"SubmissionID_1", "ParticipantID_1", "TrialID_1", "2024-11-18T12:00:00Z", "+10:00", "0.1", "0.2", "0.3"},
			{"SubmissionID_1", "ParticipantID_1", "TrialID_1", "2024-11-18T12:00:00Z", "+10:00", "0.1", "0.2", "0.3"},
			{"SubmissionID_1", "ParticipantID_1", "TrialID_1", "2024-11-18T12:00:00Z", "+10:00", "0.1", "0.2", "0.3"},
			// Add more rows as needed
		}
		for _, row := range testrows {
			ptrRow := make([]*string, len(row))
			for i, v := range row {
				ptrRow[i] = &v
			}

			if err = pw.WriteString(ptrRow); err != nil {
				log.Printf("Error writing row: %v", err)
			}
		}
	}

	// Finalize the writer
	if err = pw.WriteStop(); err != nil {
		log.Fatalf("Error finalizing Parquet file: %v", err)
	}

	log.Println("Rows successfully written to Parquet")
}
