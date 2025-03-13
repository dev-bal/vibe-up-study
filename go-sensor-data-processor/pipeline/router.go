package pipeline

import (
	"fmt"
	"sync"
	"vibe-up/sensor-data/models"
)

func Router(id int, input chan *models.Sensor, router map[string]chan *models.Sensor, nwriters int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create multiple workers for each route
	for d := range input {
		route := (*d).ByRoute()
		if _, ok := router[route]; !ok {
			router[route] = make(chan *models.Sensor, 10)
			fmt.Println("Creating new route for", route)

			// Launch multiple workers for each route
			for i := range nwriters {
				workerID := i
				wg.Add(1)
				go CSVWriter(fmt.Sprintf("%s_%d_%d", route, id, workerID), router[route], wg)
			}
		}
		router[route] <- d
	}

	for _, ch := range router {
		close(ch)
	}

}
