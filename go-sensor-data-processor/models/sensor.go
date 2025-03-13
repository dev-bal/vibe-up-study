package models

import "fmt"

type Sensor interface {
	// FromRecord(record *Record) (DataProcessor, error)
	Unnest() ([][]string, error)
	Summarise() []string
	Length() int
	ByRoute() string
}

func NewSensor(sensorType string) (func([]string) (Sensor, error), error) {
	switch sensorType {
	case "accelerometer":
		return NewAccelerometer, nil
	// case "gyroscope":
	// return NewGyroscope, nil
	// Add more sensor types here as needed
	default:
		return nil, fmt.Errorf("unknown sensor type: %s", sensorType)
	}
}
