package main

import (
	"log"
	"sync"

	"vibe-up/sensor-data/models"
	"vibe-up/sensor-data/pipeline"
)

func main() {
	config := struct {
		InputDir   string
		NumRouters int
		NumWorkers int
		NumWriters int
	}{
		InputDir:   "/input",
		NumRouters: 2,
		NumWorkers: 25,
		NumWriters: 2,
	}

	factory, err := models.NewSensor("accelerometer")
	if err != nil {
		log.Fatalf("Failed to create sensor factory: %v", err)
	}

	routerChan := make(chan *models.Sensor, 10)
	summaryChan := make(chan *models.Sensor, 10)
	fileChan := pipeline.GetFiles(config.InputDir)

	var routerWG sync.WaitGroup
	var fileWG sync.WaitGroup

	// Start file processing workers
	go pipeline.SummaryStatistics(summaryChan)

	for range config.NumWorkers {
		fileWG.Add(1)
		go func() {
			defer fileWG.Done()
			for filename := range fileChan {
				if err := pipeline.ProcessFile(factory, filename, routerChan, summaryChan); err != nil {
					log.Printf("Error processing file %s: %v", filename, err)
				}
			}
		}()
	}

	for i := range config.NumRouters {
		routerWG.Add(1)
		go func() {
			localRouter := make(map[string]chan *models.Sensor)
			pipeline.Router(i, routerChan, localRouter, config.NumWriters, &routerWG)
		}()
	}

	fileWG.Wait()
	close(routerChan)
	close(summaryChan)
	routerWG.Wait()
}
