package models

import (
	"encoding/json"
	"fmt"
	"math"
)

// Record struct to hold data for each CSV row
type Accelerometer struct {
	SubmissionID  string
	ParticipantID string
	TrialID       string
	Timestamp     string
	Timezone      string
	Payload       JSONData
}

// JSONData struct to hold the "values" JSON field
type JSONData struct {
	Values [][]float64 `json:"values"`
}

func NewAccelerometer(row []string) (Sensor, error) {
	var jsonData JSONData
	if err := json.Unmarshal([]byte(row[14]), &jsonData); err != nil {
		return nil, fmt.Errorf("JSON parse fail for SubmissionId %s: %v", row[0], err)
	}
	record := &Accelerometer{
		SubmissionID:  row[0],
		ParticipantID: row[1],
		TrialID:       row[3],
		Payload:       jsonData,
		Timestamp:     row[15],
		Timezone:      row[19],
	}
	return record, nil
}

func (r *Accelerometer) Unnest() ([][]string, error) {
	var rows [][]string
	for _, value := range r.Payload.Values {
		row := []string{
			r.SubmissionID,
			r.ParticipantID,
			r.TrialID,
			r.Timestamp,
			r.Timezone,
			fmt.Sprintf("%v", value[0]),
			fmt.Sprintf("%v", value[1]),
			fmt.Sprintf("%v", value[2]),
			fmt.Sprintf("%v", value[3]),
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (r *Accelerometer) Summarise() []string {
	var maxOffset float64
	var rootMeanSquare float64
	n := len(r.Payload.Values)
	if n > 0 {
		maxOffset = r.Payload.Values[n-1][0]
		var sumOfSquares float64
		for _, vector := range r.Payload.Values {
			var magnitudeSquared float64
			for _, component := range vector {
				magnitudeSquared += component * component
			}
			sumOfSquares += magnitudeSquared
		}
		rootMeanSquare = math.Sqrt(sumOfSquares / float64(n))
	}

	return []string{
		r.SubmissionID,
		r.ParticipantID,
		r.TrialID,
		r.Timestamp,
		r.Timezone,
		fmt.Sprintf("%v", n),
		fmt.Sprintf("%v", maxOffset),
		fmt.Sprintf("%v", rootMeanSquare),
		// MinX:          min(r.Payload.Values[0]),
		// MinY:          min(r.Payload.Values[1]),
		// MinZ:          min(r.Payload.Values[2]),
	}
}

func (r *Accelerometer) Length() int {
	return len(r.Payload.Values)
}

func (r *Accelerometer) ByRoute() string {
	return r.TrialID
}
